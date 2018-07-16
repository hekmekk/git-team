FROM golang:1.9-stretch

MAINTAINER Rea Sand <hekmek@posteo.de>

RUN apt-get update && apt-get -y install \
	libgit2-24 \
	libgit2-dev \
	man-db \
	build-essential \
	ruby \
	ruby-dev \
	rubygems

RUN gem install --no-ri --no-rdoc fpm

ENV GOPATH=/go

WORKDIR ${GOPATH}/src/github.com/hekmekk/git-team
COPY src/ .

RUN make

CMD ["fpm", "--version"]
