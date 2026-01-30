VERSION := $(shell grep -E "version\s+=" `pwd`/main.go | awk -F '"' '{print $$2}')

GOOS :=
GOARCH :=
prefix :=
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin
	prefix:=/usr/local
endif
ifeq ($(UNAME_M),arm64)
	GOARCH=arm64
endif
ifeq ($(UNAME_S),Linux)
	GOOS=linux
	prefix:=/usr
endif
ifeq ($(UNAME_M),x86_64)
	GOARCH=amd64
endif

exec_prefix := $(prefix)
bindir := $(exec_prefix)/bin
datarootdir := $(prefix)/share
man1dir := $(datarootdir)/man/man1

bash_completion_dir := $(datarootdir)/bash-completion/completions
zsh_completion_dir := $(datarootdir)/zsh/site-functions

CURR_DIR := $(shell pwd)

BATS_FILE :=
ifdef CASE
	BATS_FILE=$(CASE).bats
endif

BATS_FILTER :=
ifdef FILTER
	BATS_FILTER=--filter $(FILTER)
endif

all: build man-page completion

update-and-cleanup-deps:
	go get -u -t
	go mod tidy

deps:
	go get -t
	go mod download

test: go-test hookscript-tests

verify: test acceptance-tests

go-test: deps
	go test -cover ./src/...

fmt: deps
	go fmt main.go
	go fmt ./src/...

build: deps
ifndef GOPATH
	$(error GOPATH is not set)
endif
	mkdir -p $(CURR_DIR)/target/bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(CURR_DIR)/target/bin ./...
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page: deps
	mkdir -p $(CURR_DIR)/target/man/
	go run $(CURR_DIR)/main.go --generate-man-page > $(CURR_DIR)/target/man/git-team.1
	gzip -f $(CURR_DIR)/target/man/git-team.1

completion: bash_completion zsh_completion

bash_completion: deps
	mkdir -p $(CURR_DIR)/target/completion/bash
	go run $(CURR_DIR)/main.go completion bash > $(CURR_DIR)/target/completion/bash/git-team.bash

zsh_completion: deps
	mkdir -p $(CURR_DIR)/target/completion/zsh
	go run $(CURR_DIR)/main.go completion zsh > $(CURR_DIR)/target/completion/zsh/git-team.zsh

install:
	@echo "[INFO] Installing into $(bindir)/ ..."
	mkdir -p $(bindir)
	install $(CURR_DIR)/target/bin/git-team $(bindir)/git-team
	mkdir -p $(man1dir)
	install -m "0644" $(CURR_DIR)/target/man/git-team.1.gz $(man1dir)/git-team.1.gz
	mkdir -p $(bash_completion_dir)
	install -m "0644" $(CURR_DIR)/target/completion/bash/git-team.bash $(bash_completion_dir)/git-team
	echo "[INFO] bash: Don't forget to source $(bash_completion_dir)/git-team"
	mkdir -p $(zsh_completion_dir)
	install -m "0644" $(CURR_DIR)/target/completion/zsh/git-team.zsh $(zsh_completion_dir)/_git-team
	echo "[INFO] zsh: Don't forget to source $(zsh_completion_dir)/_git-team"

uninstall:
	rm -f $(bindir)/git-team
	rm -f $(bash_completion_dir)/git-team
	rm -f $(zsh_completion_dir)/_git-team
	rm -f $(man1dir)/git-team.1.gz

clean:
	rm -rf $(CURR_DIR)/target

.PHONY: acceptance-tests
acceptance-tests:
	docker build -t git-team-acceptance-tests . -f acceptance-tests.Dockerfile
	docker run --rm -v $(CURR_DIR)/acceptance-tests:/acceptance-tests:ro git-team-acceptance-tests --formatter tap /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)

.PHONY: hookscript-tests
hookscript-tests:
	docker build -t git-team-hookscript-tests . -f hookscript-tests.Dockerfile
	docker run --rm -v $(CURR_DIR)/hookscript-tests:/hookscript-tests:ro git-team-hookscript-tests --formatter tap /hookscript-tests/

