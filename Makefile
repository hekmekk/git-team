VERSION:=$(shell grep "version =" `pwd`/cmd/git-team/main.go | awk -F '"' '{print $$2}' | cut -c2-)

CURR_DIR:=$(shell pwd)

UNAME_S:= $(shell uname -s)
BASH_COMPLETION_PREFIX:=
GOOS:=linux
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin
	BASH_COMPLETION_PREFIX:=/usr/local
endif
HOOKS_DIR:=/usr/local/etc/git-team/hooks

BATS_FILE:=
ifdef CASE
	BATS_FILE=$(CASE).bats
endif

BATS_FILTER:=
ifdef FILTER
	BATS_FILTER=--filter $(FILTER)
endif

all: deps fmt build man-page

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
	mkdir -p $(CURR_DIR)/pkg/target/bin
	mv $(GOPATH)/bin/git-team $(CURR_DIR)/pkg/target/bin/git-team
	mv $(GOPATH)/bin/prepare-commit-msg $(CURR_DIR)/pkg/target/bin/prepare-commit-msg
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page:
	mkdir -p $(CURR_DIR)/pkg/target/man/
	go run $(CURR_DIR)/cmd/git-team/main.go --help-man > $(CURR_DIR)/pkg/target/man/git-team.1
	gzip -f $(CURR_DIR)/pkg/target/man/git-team.1

install:
	@echo "[INFO] Installing into $(BIN_PREFIX)/bin/ ..."
	install $(CURR_DIR)/pkg/target/bin/git-team $(BIN_PREFIX)/bin/git-team
	mkdir -p $(HOOKS_DIR)
	install $(CURR_DIR)/pkg/target/bin/prepare-commit-msg $(HOOKS_DIR)/prepare-commit-msg
	install $(CURR_DIR)/git-hooks/proxy.sh /usr/local/etc/git-team/hooks/proxy.sh
	$(CURR_DIR)/git-hooks/install_symlinks.sh
	mkdir -p /usr/local/share/man/man1
	install -m "0644" pkg/target/man/git-team.1.gz /usr/local/share/man/man1/git-team.1.gz
	install -m "0644" bash_completion/git-team.bash $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source $(BASH_COMPLETION_PREFIX)/etc/bash_completion"

uninstall:
	rm -f $(BIN_PREFIX)/bin/git-team
	rm -f /etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz
	rm -rf $(HOOKS_DIR)

package-build: clean
	mkdir -p pkg/src/
	cp Makefile pkg/src/
	cp go.mod pkg/src/
	cp -r cmd pkg/src/
	cp -r src pkg/src/
	cp -r bash_completion pkg/src/
	cp -r git-hooks pkg/src/
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-pkg:v$(VERSION) pkg/

deb rpm: package-build
	mkdir -p pkg/target/$@
	chown -R $(shell id -u):$(shell id -g) pkg/target/$@
	docker run --rm -h git-team-pkg -v $(CURR_DIR)/pkg/target/$@:/pkg-target git-team-pkg:v$(VERSION) fpm \
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
		pkg/target/bin/git-team=$(BIN_PREFIX)/bin/git-team \
		pkg/target/bin/prepare-commit-msg=$(HOOKS_DIR)/prepare-commit-msg \
		git-hooks/proxy.sh=$(HOOKS_DIR)/proxy.sh \
		bash_completion/git-team.bash=/etc/bash_completion.d/git-team \
		pkg/target/man/git-team.1.gz=/usr/share/man/man1/git-team.1.gz

package: deb rpm

release-github: package
ifndef GITHUB_API_TOKEN
	$(error GITHUB_API_TOKEN is not set)
endif
	lua scripts/release-github.lua \
		--github-api-token $(GITHUB_API_TOKEN) \
		--git-team-version v$(VERSION) \
		--git-team-deb-path $(CURR_DIR)/pkg/target/deb/git-team_$(VERSION)_amd64.deb

release: release-github

clean:
	rm -f $(CURR_DIR)/git-team
	rm -rf $(CURR_DIR)/pkg/src/
	rm -rf $(CURR_DIR)/pkg/target/
	rm -rf $(CURR_DIR)/acceptance-tests/src/
	rm -rf $(CURR_DIR)/acceptance-tests/git-hooks/

purge: clean uninstall
	git config --global --unset-all commit.template
	git config --global --unset-all core.hooksPath

docker-build: clean
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .
	docker tag git-team-run:v$(VERSION) git-team-run:latest

.PHONY: acceptance-tests
acceptance-tests:
	rm -rf $(CURR_DIR)/acceptance-tests/src
	mkdir -p $(CURR_DIR)/acceptance-tests/src/
	cp $(CURR_DIR)/go.* $(CURR_DIR)/acceptance-tests/src/
	cp -r $(CURR_DIR)/cmd $(CURR_DIR)/acceptance-tests/src/
	cp -r $(CURR_DIR)/src $(CURR_DIR)/acceptance-tests/src/
	cp -r $(CURR_DIR)/git-hooks $(CURR_DIR)/acceptance-tests/git-hooks
	docker build -t git-team-acceptance-tests $(CURR_DIR)/acceptance-tests
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/acceptance-tests/cases:/acceptance-tests git-team-acceptance-tests --pretty /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)
