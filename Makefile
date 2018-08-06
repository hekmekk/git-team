VERSION:=0.1.0

all: test fmt build man-page

deps:
	go get

test: deps
	go test -short git-team.go
	go test -short ./core/...

fmt:
	go fmt

build: deps
	go build

man-page:
	mkdir -p man/
	go run git-team.go --help-man > man/git-team.1
	gzip -f man/git-team.1

install:
	install git-team /usr/bin/git-team
	install --mode="0644" man/git-team.1.gz /usr/share/man/man1/git-team.1.gz
	install --mode="0644" bash_completion/git-team.bash /etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source /etc/bash_completion"

package-build:
	mkdir -p pkg/src/
	cp Makefile pkg/src/
	cp git-team.go pkg/src/
	cp -r core pkg/src/
	cp -r bash_completion pkg/src/
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-pkg:v$(VERSION) pkg/

package: package-build
	mkdir -p pkg/target/
	chown -R $(shell id -u):$(shell id -g) pkg/target/
	docker run --rm -v `pwd`/pkg/target:/target git-team-pkg:v$(VERSION) fpm \
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
		-p /target \
		git-team=/usr/bin/git-team \
		bash_completion/git-team.bash=/etc/bash_completion.d/git-team \
		man/git-team.1.gz=/usr/share/man/man1/git-team.1.gz

release:
	@echo "nope... :D"

clean:
	rm -f git-team
	rm -rf man/
	rm -rf pkg/src/
	rm -rf pkg/target/

purge: clean
	rm -f /usr/bin/git-team
	rm -f /etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz
	git config --global --remove-section team.alias || true
	git config --global --remove-section commit || true
	git config --remove-section team.alias || true
	git config --remove-section commit || true

docker-build:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg USERNAME=$(USER) -t git-team-run:v$(VERSION) .

docker-run: docker-build
	mkdir -p /home/$(USER)/.config/git-team
	chown -R $(shell id -u):$(shell id -g) /home/$(USER)/.config/git-team
	docker run --rm -v /home/$(USER)/:/home/$(USER)/ git-team-run:v$(VERSION) git team --help
