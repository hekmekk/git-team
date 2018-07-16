FROM golang:1.9-stretch

MAINTAINER Rea Sand <hekmek@posteo.de>

ARG USERNAME=git-team-run
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USERNAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USERNAME

RUN apt-get update && apt-get -y install \
	libgit2-24 \
	libgit2-dev \
	man-db

WORKDIR /go/src/github.com/hekmekk/git-team

ENV GOPATH=/go

COPY . .

RUN make

RUN make install

USER $USERNAME

CMD ["/usr/bin/git-team", "--help"]
