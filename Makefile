VERSION:=$(shell grep "version =" `pwd`/cmd/git-team/main.go | awk -F '"' '{print $$2}' | cut -c2-)

CURR_DIR:=$(shell pwd)

UNAME_S:= $(shell uname -s)
BASH_COMPLETION_PREFIX:=
ifeq ($(UNAME_S),Darwin)
	BASH_COMPLETION_PREFIX:=/usr/local
endif

all: test fmt build man-page

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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./cmd/...
	mkdir -p $(CURR_DIR)/pkg/target/bin
	mv $(GOPATH)/bin/git-team $(CURR_DIR)/pkg/target/bin/git-team
	mv $(GOPATH)/bin/prepare-commit-msg $(CURR_DIR)/pkg/target/bin/prepare-commit-msg
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page:
	mkdir -p $(CURR_DIR)/pkg/target/man/
	go run $(CURR_DIR)/cmd/git-team/main.go --help-man > pkg/target/man/git-team.1
	gzip -f $(CURR_DIR)/pkg/target/man/git-team.1

install:
	install $(CURR_DIR)/pkg/target/bin/git-team /usr/local/bin/git-team
	mkdir -p /usr/local/share/.config/git-team/hooks
	install $(CURR_DIR)/pkg/target/bin/prepare-commit-msg /usr/local/share/.config/git-team/hooks/prepare-commit-msg
	mkdir -p /usr/local/share/man/man1
	install -m "0644" pkg/target/man/git-team.1.gz /usr/local/share/man/man1/git-team.1.gz
	install -m "0644" bash_completion/git-team.bash $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source $(BASH_COMPLETION_PREFIX)/etc/bash_completion"

uninstall:
	rm -f /usr/local/bin/git-team
	rm -f /etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz
	rm -rf /usr/local/share/.config/git-team

package-build: clean
	mkdir -p pkg/src/
	cp Makefile pkg/src/
	cp go.mod pkg/src/
	cp -r cmd pkg/src/
	cp -r src pkg/src/
	cp -r bash_completion pkg/src/
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-pkg:v$(VERSION) pkg/

package: package-build
	mkdir -p pkg/target/deb
	chown -R $(shell id -u):$(shell id -g) pkg/target/deb
	docker run --rm -h git-team-pkg -v $(CURR_DIR)/pkg/target/deb:/deb-target git-team-pkg:v$(VERSION) fpm \
		-f \
		-s dir \
		-t deb \
		-n "git-team" \
		-v $(VERSION) \
		-m "Rea Sand <hekmek@posteo.de>" \
		--url "https://github.com/hekmekk/git-team" \
		--license "MIT" \
		--description "git-team - commit template provisioning with co-authors" \
		--deb-no-default-config-files \
		-p /deb-target \
		pkg/target/bin/git-team=/usr/bin/git-team \
		pkg/target/bin/prepare-commit-msg=/usr/local/share/.config/git-team/hooks/prepare-commit-msg \
		bash_completion/git-team.bash=/etc/bash_completion.d/git-team \
		pkg/target/man/git-team.1.gz=/usr/share/man/man1/git-team.1.gz

release:
# TODO: requires 'package' target
ifndef GITHUB_API_TOKEN
	$(error GITHUB_API_TOKEN is not set)
endif
	lua scripts/release.lua \
		--github-api-token $(GITHUB_API_TOKEN) \
		--git-team-version v$(VERSION) \
		--git-team-deb-path $(CURR_DIR)/pkg/target/deb/git-team_$(VERSION)_amd64.deb

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
	docker build -t git-team-acceptance-tests $(CURR_DIR)/acceptance-tests
	docker run --rm -v $(CURR_DIR)/acceptance-tests/cases:/acceptance-tests git-team-acceptance-tests --tap /acceptance-tests

