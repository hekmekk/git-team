VERSION:=$(shell grep -E "version\s+=" `pwd`/cmd/git-team/main.go | awk -F '"' '{print $$2}')

CURR_DIR:=$(shell pwd)

UNAME_S:= $(shell uname -s)
BASH_COMPLETION_PREFIX:=
GOOS:=linux
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin
	BASH_COMPLETION_PREFIX:=/usr/local
endif
# Note: this is currently still hard-coded
HOOKS_DIR:=/usr/local/etc/git-team/hooks

BATS_FILE:=
ifdef CASE
	BATS_FILE=$(CASE).bats
endif

BATS_FILTER:=
ifdef FILTER
	BATS_FILTER=--filter $(FILTER)
endif

all: fmt build man-page process-hook-templates

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

mo:
	curl -sSL https://git.io/get-mo -o mo
	chmod +x mo

process-hook-templates: mo
	hooks_dir=$(HOOKS_DIR) ./mo $(CURR_DIR)/git-hooks/prepare-commit-msg.sh.mo > $(CURR_DIR)/git-hooks/prepare-commit-msg.sh
	hooks_dir=$(HOOKS_DIR) ./mo $(CURR_DIR)/git-hooks/install_symlinks.sh.mo > $(CURR_DIR)/git-hooks/install_symlinks.sh
	chmod +x $(CURR_DIR)/git-hooks/install_symlinks.sh

install:
	@echo "[INFO] Installing into $(BIN_PREFIX)/bin/ ..."
	install $(CURR_DIR)/target/bin/git-team $(BIN_PREFIX)/bin/git-team
	mkdir -p $(HOOKS_DIR)
	install $(CURR_DIR)/target/bin/prepare-commit-msg-git-team $(HOOKS_DIR)/prepare-commit-msg-git-team
	install $(CURR_DIR)/git-hooks/proxy.sh $(HOOKS_DIR)/proxy.sh
	install $(CURR_DIR)/git-hooks/prepare-commit-msg.sh $(HOOKS_DIR)/prepare-commit-msg
	$(CURR_DIR)/git-hooks/install_symlinks.sh
	mkdir -p /usr/local/share/man/man1
	install -m "0644" $(CURR_DIR)/target/man/git-team.1.gz /usr/local/share/man/man1/git-team.1.gz
	@if [ -d "$(BASH_COMPLETION_PREFIX)/etc/bash_completion.d" ]; then \
		install -m "0644" $(CURR_DIR)/bash_completion/git-team.bash $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/git-team; \
		echo "[INFO] Don't forget to source $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/*"; \
	fi

uninstall:
	rm -f $(BIN_PREFIX)/bin/git-team
	rm -f $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz
	rm -rf $(HOOKS_DIR)

export-signing-key: clean
ifndef GPG_SIGNING_KEY_ID
	$(error GPG_SIGNING_KEY_ID is not set)
endif
	gpg --armor --export-secret-keys $(GPG_SIGNING_KEY_ID) > $(CURR_DIR)/signing-key.asc

package-build: export-signing-key
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-pkg:v$(VERSION) . -f pkg.Dockerfile

deb rpm: clean package-build process-hook-templates
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
		--after-install git-hooks/install_symlinks.sh \
		--deb-no-default-config-files \
		--rpm-sign \
		-p /pkg-target \
		target/bin/git-team=$(BIN_PREFIX)/bin/git-team \
		target/bin/prepare-commit-msg-git-team=$(HOOKS_DIR)/prepare-commit-msg-git-team \
		git-hooks/proxy.sh=$(HOOKS_DIR)/proxy.sh \
		git-hooks/prepare-commit-msg.sh=$(HOOKS_DIR)/prepare-commit-msg \
		bash_completion/git-team.bash=/etc/bash_completion.d/git-team \
		target/man/git-team.1.gz=/usr/share/man/man1/git-team.1.gz

show-checksums: package-build
	find $(CURR_DIR)/target/ -type f -exec sha256sum {} \;

package: rpm deb show-checksums

clean:
	rm -f $(CURR_DIR)/git-team
	rm -f $(CURR_DIR)/signing-key.asc
	rm -f $(CURR_DIR)/git-hooks/prepare-commit-msg.sh
	rm -f $(CURR_DIR)/git-hooks/install_symlinks.sh
	rm -rf $(CURR_DIR)/target
	rm -rf $(CURR_DIR)/acceptance-tests/src/
	rm -rf $(CURR_DIR)/acceptance-tests/git-hooks/

purge: clean uninstall
	git config --global --unset-all commit.template
	git config --global --unset-all core.hooksPath

docker-build: clean
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .
	docker tag git-team-run:v$(VERSION) git-team-run:latest

docker-run: docker-build
	docker run git-team-run:v$(VERSION) --help

.PHONY: acceptance-tests
acceptance-tests:
	docker build -t git-team-acceptance-tests . -f acceptance-tests.Dockerfile
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/acceptance-tests:/acceptance-tests git-team-acceptance-tests --pretty /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)
