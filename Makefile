VERSION:=0.0.1-alpha1

all: test fmt build man_page

os_deps:
	apt-get -y install \
		libgit2-24 \
		libgit2-dev \
		docker-ce

deps:
	go get

test: deps
	go test -short ./...

fmt:
	go fmt

build: deps
	go build

man_page:
	mkdir -p man/
	go run git-team.go --help-man > man/git-team.1
	gzip -f man/git-team.1

install:
	install git-team /usr/bin/git-team
	install --mode="0644" man/git-team.1.gz /usr/share/man/man1/git-team.1.gz
	install --mode="0644" bash_completion/git-team.bash /etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source /etc/bash_completion"

package_build:
	mkdir -p pkg/src/
	cp Makefile pkg/src/
	cp git-team.go pkg/src/
	cp -r core pkg/src/
	cp -r bash_completion pkg/src/
	docker build -t git-team-pkg:$(VERSION) pkg/

package: package_build
	docker run --rm -v `pwd`:/src -v `pwd`/pkg/target:/target git-team-pkg:$(VERSION) fpm \
		-f \
		-s dir \
		-t deb \
		-n "git-team" \
		-v $(VERSION) \
		--url "https://github.com/hekmekk/git-team" \
		--license "MIT" \
		--description "git-team - commit template provisioning with co-authors" \
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

docker_build:
	 docker build -t git-team-docker:$(VERSION) .

docker: docker_build
	 docker run --rm git-team-docker:$(VERSION) git team --help
