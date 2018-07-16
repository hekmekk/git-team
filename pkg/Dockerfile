FROM golang:1.9-stretch

MAINTAINER Rea Sand <hekmek@posteo.de>

ARG USERNAME=git-team-pkg
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USERNAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USERNAME

RUN apt-get update && apt-get -y install \
	libgit2-24 \
	libgit2-dev \
	man-db \
	build-essential \
	ruby \
	ruby-dev \
	rubygems

RUN gem install --no-ri --no-rdoc fpm

WORKDIR /go/src/github.com/hekmekk/git-team

RUN mkdir -p /target
RUN chown -R $USERNAME:$USERNAME /target
RUN chown -R $USERNAME:$USERNAME /go

USER $USERNAME

ENV GOPATH=/go

COPY src/ .

RUN make

CMD ["fpm", "--version"]
