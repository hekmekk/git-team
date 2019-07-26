VERSION:=1.1.1

UNAME_S:= $(shell uname -s)
BASH_COMPLETION_PREFIX:=
ifeq ($(UNAME_S),Darwin)
	BASH_COMPLETION_PREFIX:=/usr/local
endif

all: test fmt build man-page

tidy:
	go mod tidy

deps: tidy
	go get

test: deps
	go test -cover .
	go test -cover ./src/...

fmt:
	go fmt

build: deps
	go build -o pkg/target/bin/git-team

man-page:
	mkdir -p pkg/target/man/
	go run main.go --help-man > pkg/target/man/git-team.1
	gzip -f pkg/target/man/git-team.1

install:
	install pkg/target/bin/git-team /usr/local/bin/git-team
	mkdir -p /usr/local/share/man/man1
	install -m "0644" pkg/target/man/git-team.1.gz /usr/local/share/man/man1/git-team.1.gz
	install -m "0644" bash_completion/git-team.bash $(BASH_COMPLETION_PREFIX)/etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source $(BASH_COMPLETION_PREFIX)/etc/bash_completion"

uninstall:
	rm -f /usr/bin/git-team
	rm -f /etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz

package-build: clean
	mkdir -p pkg/src/
	cp Makefile pkg/src/
	cp main.go pkg/src/
	cp go.mod pkg/src/
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
		bash_completion/git-team.bash=/etc/bash_completion.d/git-team \
		pkg/target/man/git-team.1.gz=/usr/share/man/man1/git-team.1.gz

release:
	@echo "nope... :D"

clean:
	rm -f git-team
	rm -rf pkg/src/
	rm -rf pkg/target/

purge: clean uninstall
	git config --global --remove-section team.alias || true
	git config --global --remove-section commit || true
	git config --remove-section team.alias || true
	git config --remove-section commit || true

docker-build: clean
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .

docker-run: docker-build
	mkdir -p /home/$(USER)/.config/git-team
	chown -R $(shell id -u):$(shell id -g) /home/$(USER)/.config/git-team
	docker run --rm -h git-team-run -v /home/$(USER)/:/home/$(USER)/ git-team-run:v$(VERSION) git team --help
