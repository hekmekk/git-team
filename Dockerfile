FROM golang:1.9-stretch

MAINTAINER Rea Sand <hekmek@posteo.de>

RUN apt-get update && apt-get -y install libgit2-24 libgit2-dev man-db

ENV GOPATH=/go

WORKDIR ${GOPATH}/src/github.com/hekmekk/git-team
COPY . .

RUN make

RUN make install

CMD ["/usr/bin/git-team", "--help"]
