VERSION := $(shell grep -E "version\s+=" `pwd`/main.go | awk -F '"' '{print $$2}')

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

bash_completion_dir := $(datarootdir)/bash-completion/completions

CURR_DIR := $(shell pwd)

BATS_FILE :=
ifdef CASE
	BATS_FILE=$(CASE).bats
endif

BATS_FILTER :=
ifdef FILTER
	BATS_FILTER=--filter $(FILTER)
endif

all: fmt build man-page completion

deps:
	go mod download

test: clean go-test hookscript-tests

verify: test acceptance-tests

mocks:
	docker run --rm --user "$(shell id -u):$(shell id -g)" -v "$(CURR_DIR):/src" -w /src vektra/mockery:v2.8 --dir=src/ --all --keeptree

go-test: mocks deps
	go test -cover ./src/...

fmt: deps
	go fmt main.go
	go fmt ./src/...

build: clean deps
ifndef GOPATH
	$(error GOPATH is not set)
endif
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=amd64 go install ./...
	mkdir -p $(CURR_DIR)/target/bin
	mv $(GOPATH)/bin/git-team $(CURR_DIR)/target/bin/git-team
	@echo "[INFO] Successfully built git-team version v$(VERSION)"

man-page: clean deps
	mkdir -p $(CURR_DIR)/target/man/
	go run $(CURR_DIR)/main.go --generate-man-page > $(CURR_DIR)/target/man/git-team.1
	gzip -f $(CURR_DIR)/target/man/git-team.1

completion: clean deps
	mkdir -p $(CURR_DIR)/target/completion/bash
	go run $(CURR_DIR)/main.go completion bash > $(CURR_DIR)/target/completion/bash/git-team.bash

install:
	@echo "[INFO] Installing into $(bindir)/ ..."
	mkdir -p $(bindir)
	install $(CURR_DIR)/target/bin/git-team $(bindir)/git-team
	mkdir -p $(man1dir)
	install -m "0644" $(CURR_DIR)/target/man/git-team.1.gz $(man1dir)/git-team.1.gz
	mkdir -p $(bash_completion_dir)
	install -m "0644" $(CURR_DIR)/target/completion/bash/git-team.bash $(bash_completion_dir)/git-team
	echo "[INFO] Don't forget to source $(bash_completion_dir)/git-team"

uninstall:
	rm -f $(bindir)/git-team
	rm -f $(bash_completion_dir)/git-team
	rm -f $(man1dir)/git-team.1.gz

clean:
	rm -rf $(CURR_DIR)/mocks
	rm -rf $(CURR_DIR)/target

.PHONY: acceptance-tests
acceptance-tests:
	docker build -t git-team-acceptance-tests . -f acceptance-tests.Dockerfile
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/acceptance-tests:/acceptance-tests git-team-acceptance-tests --pretty /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)

.PHONY: hookscript-tests
hookscript-tests:
	docker build -t git-team-hookscript-tests . -f hookscript-tests.Dockerfile
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/hookscript-tests:/hookscript-tests git-team-hookscript-tests --pretty /hookscript-tests/

