VERSION:=$(shell grep "version =" `pwd`/cmd/git-team/main.go | awk -F '"' '{print $$2}' | cut -c2-)

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
	mkdir -p $(shell pwd)/pkg/target/bin
	mv $(GOPATH)/bin/git-team $(shell pwd)/pkg/target/bin/git-team
	mv $(GOPATH)/bin/prepare-commit-msg $(shell pwd)/pkg/target/bin/prepare-commit-msg
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page:
	mkdir -p $(shell pwd)/pkg/target/man/
	go run $(shell pwd)/cmd/git-team/main.go --help-man > pkg/target/man/git-team.1
	gzip -f $(shell pwd)/pkg/target/man/git-team.1

install:
	install $(shell pwd)/pkg/target/bin/git-team /usr/local/bin/git-team
	mkdir -p /usr/local/share/.config/git-team/hooks
	install $(shell pwd)/pkg/target/bin/prepare-commit-msg /usr/local/share/.config/git-team/hooks/prepare-commit-msg
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
	docker run --rm -h git-team-pkg -v $(shell pwd)/pkg/target/deb:/deb-target git-team-pkg:v$(VERSION) fpm \
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
	@echo "nope... :D"

clean:
	rm -f git-team
	rm -rf pkg/src/
	rm -rf pkg/target/
	rm -rf acceptance-tests/src

purge: clean uninstall
	git config --global --unset-all commit.template
	git config --global --unset-all core.hooksPath

docker-build: clean
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .
	docker tag git-team-run:v$(VERSION) git-team-run:latest

.PHONY: acceptance-tests
acceptance-tests: clean
	mkdir -p acceptance-tests/src/
	cp go.* acceptance-tests/src/
	cp -r cmd acceptance-tests/src/
	cp -r src acceptance-tests/src/
	docker build -t git-team-acceptance-tests $(shell pwd)/acceptance-tests
	docker run --rm -v $(shell pwd)/acceptance-tests/cases:/acceptance-tests git-team-acceptance-tests --tap /acceptance-tests

