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
AVAILABLE_GIT_HOOKS:=applypatch-msg commit-msg fsmonitor-watchman p4-pre-submit post-applypatch post-checkout post-commit post-index-change post-merge post-receive post-rewrite post-update pre-applypatch pre-auto-gc pre-commit pre-push pre-rebase pre-receive push-to-checkout sendemail-validate update

all: deps build man-page

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
	for hook in $(AVAILABLE_GIT_HOOKS) ; do \
		install $(CURR_DIR)/git-hooks/proxy.sh /usr/local/etc/git-team/hooks/$$hook ; \
	done
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
		-m "Rea Sand <hekmek@posteo.de>" \
		--url "https://github.com/hekmekk/git-team" \
		--license "MIT" \
		--description "git-team - commit message enhancement with co-authors" \
		--deb-no-default-config-files \
		-p /pkg-target \
		pkg/target/bin/git-team=$(BIN_PREFIX)/bin/git-team \
		pkg/target/bin/prepare-commit-msg=$(HOOKS_DIR)/prepare-commit-msg \
		git-hooks/proxy.sh=$(HOOKS_DIR)/applypatch-msg \
		git-hooks/proxy.sh=$(HOOKS_DIR)/commit-msg \
		git-hooks/proxy.sh=$(HOOKS_DIR)/fsmonitor-watchman \
		git-hooks/proxy.sh=$(HOOKS_DIR)/p4-pre-submit \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-applypatch \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-checkout \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-commit \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-index-change \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-merge \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-receive \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-rewrite \
		git-hooks/proxy.sh=$(HOOKS_DIR)/post-update \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-applypatch \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-auto-gc \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-commit \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-push \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-rebase \
		git-hooks/proxy.sh=$(HOOKS_DIR)/pre-receive \
		git-hooks/proxy.sh=$(HOOKS_DIR)/push-to-checkout \
		git-hooks/proxy.sh=$(HOOKS_DIR)/sendemail-validate \
		git-hooks/proxy.sh=$(HOOKS_DIR)/update \
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
	docker run --rm -v $(CURR_DIR)/acceptance-tests/cases:/acceptance-tests git-team-acceptance-tests --tap /acceptance-tests

