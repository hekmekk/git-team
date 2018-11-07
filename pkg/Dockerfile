FROM golang:1.9-stretch

LABEL maintainer Rea Sand <hekmek@posteo.de>

ARG USERNAME=git-team-pkg
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USERNAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USERNAME

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install \
	man-db \
	build-essential \
	ruby \
	ruby-dev \
	rubygems

RUN gem install --no-ri --no-rdoc fpm

WORKDIR /go/src/github.com/hekmekk/git-team

RUN mkdir -p /target
RUN chown -R $UID:$GID /target
RUN chown -R $UID:$GID /go

USER $USERNAME

ENV GOPATH=/go

COPY src/ .

RUN make

CMD ["fpm", "--version"]
