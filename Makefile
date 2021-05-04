VERSION := $(shell grep -E "version\s+=" `pwd`/cmd/git-team/main.go | awk -F '"' '{print $$2}')

GOOS :=
prefix :=
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin
	prefix:=/usr/local
endif
ifeq ($(UNAME_S),Linux)
	GOOS=linux
	prefix:=/usr
endif

exec_prefix := $(prefix)
bindir := $(exec_prefix)/bin
datarootdir := $(prefix)/share
man1dir := $(datarootdir)/man/man1
sysconfdir := $(prefix)/etc

CURR_DIR := $(shell pwd)

BATS_FILE :=
ifdef CASE
	BATS_FILE=$(CASE).bats
endif

BATS_FILTER :=
ifdef FILTER
	BATS_FILTER=--filter $(FILTER)
endif

all: fmt build man-page

tidy:
	go mod tidy

deps: tidy
	go mod download

test: deps
	go test -cover ./cmd/...
	go test -cover ./src/...

fmt:
	go fmt ./cmd/...
	go fmt ./src/...

build: clean deps
ifndef GOPATH
	$(error GOPATH is not set)
endif
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go install ./cmd/...
	mkdir -p $(CURR_DIR)/target/bin
	mv $(GOPATH)/bin/git-team $(CURR_DIR)/target/bin/git-team
	mv $(GOPATH)/bin/prepare-commit-msg $(CURR_DIR)/target/bin/prepare-commit-msg-git-team
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page: clean deps
	mkdir -p $(CURR_DIR)/target/man/
	go run $(CURR_DIR)/cmd/git-team/main.go --generate-man-page > $(CURR_DIR)/target/man/git-team.1
	gzip -f $(CURR_DIR)/target/man/git-team.1

install:
	@echo "[INFO] Installing into $(bindir)/ ..."
	mkdir -p $(bindir)
	install $(CURR_DIR)/target/bin/git-team $(bindir)/git-team
	install $(CURR_DIR)/target/bin/prepare-commit-msg-git-team $(bindir)/prepare-commit-msg-git-team
	mkdir -p $(man1dir)
	install -m "0644" $(CURR_DIR)/target/man/git-team.1.gz $(man1dir)/git-team.1.gz
	@if [ -d "$(sysconfdir)/bash_completion.d" ]; then \
		install -m "0644" $(CURR_DIR)/bash_completion/git-team.bash $(sysconfdir)/bash_completion.d/git-team; \
		echo "[INFO] Don't forget to source $(sysconfdir)/bash_completion.d/*"; \
	fi

uninstall:
	rm -f $(bindir)/git-team
	rm -f $(bindir)/prepare-commit-msg-git-team
	rm -f $(sysconfdir)/bash_completion.d/git-team
	rm -f $(man1dir)/git-team.1.gz

export-signing-key: clean
ifndef GPG_SIGNING_KEY_ID
	$(error GPG_SIGNING_KEY_ID is not set)
endif
	gpg --armor --export-secret-keys $(GPG_SIGNING_KEY_ID) > $(CURR_DIR)/signing-key.asc

package-build: export-signing-key
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-pkg:v$(VERSION) . -f pkg.Dockerfile

deb rpm: clean package-build
	mkdir -p target/$@
	chown -R $(shell id -u):$(shell id -g) target/$@
	docker run --rm -h git-team-pkg -v $(CURR_DIR)/target/$@:/pkg-target git-team-pkg:v$(VERSION) fpm \
		-f \
		-s dir \
		-t $@ \
		-n "git-team" \
		-v $(VERSION) \
		-m "git-team authors" \
		--url "https://github.com/hekmekk/git-team" \
		--architecture "x86_64" \
		--license "MIT" \
		--vendor "git-team authors" \
		--description "git-team - commit message enhancement with co-authors" \
		--depends "git" \
		--deb-no-default-config-files \
		--rpm-sign \
		-p /pkg-target \
		target/bin/git-team=$(bindir)/git-team \
		target/bin/prepare-commit-msg-git-team=$(bindir)/prepare-commit-msg-git-team \
		bash_completion/git-team.bash=$(sysconfdir)/bash_completion.d/git-team \
		target/man/git-team.1.gz=$(man1dir)/git-team.1.gz

show-checksums: package-build
	find $(CURR_DIR)/target/ -type f -exec sha256sum {} \;

package: rpm deb show-checksums

clean:
	rm -f $(CURR_DIR)/git-team
	rm -f $(CURR_DIR)/signing-key.asc
	rm -rf $(CURR_DIR)/target
	rm -rf $(CURR_DIR)/acceptance-tests/src/
	rm -rf $(CURR_DIR)/acceptance-tests/git-hooks/

docker-build: clean
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .
	docker tag git-team-run:v$(VERSION) git-team-run:latest

docker-run: docker-build
	docker run git-team-run:v$(VERSION) --help

.PHONY: acceptance-tests
acceptance-tests:
	docker build -t git-team-acceptance-tests . -f acceptance-tests.Dockerfile
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/acceptance-tests:/acceptance-tests git-team-acceptance-tests --pretty /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)
