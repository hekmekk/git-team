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

all: fmt build man-page

tidy:
	go mod tidy

deps: tidy
	go mod download

test: deps
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

install:
	@echo "[INFO] Installing into $(bindir)/ ..."
	mkdir -p $(bindir)
	install $(CURR_DIR)/target/bin/git-team $(bindir)/git-team
	mkdir -p $(man1dir)
	install -m "0644" $(CURR_DIR)/target/man/git-team.1.gz $(man1dir)/git-team.1.gz
	mkdir -p $(bash_completion_dir)
	install -m "0644" $(CURR_DIR)/bash_completion/git-team.bash $(bash_completion_dir)/git-team
	echo "[INFO] Don't forget to source $(bash_completion_dir)/git-team"

uninstall:
	rm -f $(bindir)/git-team
	rm -f $(bash_completion_dir)/git-team
	rm -f $(man1dir)/git-team.1.gz

clean:
	rm -f $(CURR_DIR)/git-team
	rm -rf $(CURR_DIR)/target
	rm -rf $(CURR_DIR)/acceptance-tests/src/
	rm -rf $(CURR_DIR)/acceptance-tests/git-hooks/

.PHONY: acceptance-tests
acceptance-tests:
	docker build -t git-team-acceptance-tests . -f acceptance-tests.Dockerfile
	docker run -e "TERM=$(TERM)" --rm -v $(CURR_DIR)/acceptance-tests:/acceptance-tests git-team-acceptance-tests --pretty /acceptance-tests/$(BATS_FILE) $(BATS_FILTER)
