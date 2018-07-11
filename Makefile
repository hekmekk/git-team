VERSION:=v0.0.1-alpha1

all: test fmt build man_page

os_deps:
	apt-get install libgit2-24 libgit2-dev docker-ce

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

install:
	install git-team /usr/bin/git-team
	install --mode="0644" man/git-team.1 /usr/share/man/man1/git-team.1
	gzip -f /usr/share/man/man1/git-team.1
	install --mode="0644" bash_completion/git-team.bash /etc/bash_completion.d/git-team
	@echo "[INFO] Don't forget to source /etc/bash_completion"

package:
	@echo "not yet :<"

release:
	@echo "nope... :D"

clean:
	rm -f git-team
	rm -rf man/
	rm -f /tmp/.git-team/STATE
	git config --global --unset commit.template || true
	git config --global --remove-section team.alias || true
	git config --remove-section team.alias || true

purge: clean
	rm -f /usr/bin/git-team
	rm -f /etc/bash_completion.d/git-team
	rm -f /usr/share/man/man1/git-team.1.gz
	rm -f $(GOPATH)/pkg/linux_amd64/gopkg.in/alecthomas/kingpin.v2.a
	rm -rf $(GOPATH)/src/gopkg.in/alecthomas/kingpin.v2
	rm -rf /tmp/.git-team

docker_build:
	 docker build -t git-team-docker:$(VERSION) .

docker: docker_build
	 docker run --rm -v /tmp/.git-team:/tmp/.git-team git-team-docker:$(VERSION) git team --help
