FROM golang:1.12-stretch

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

WORKDIR /src
COPY src/ .
RUN chown -R $UID:$GID .

RUN mkdir -p /go && chown -R $UID:$GID /go
RUN mkdir -p /deb-target && chown -R $UID:$GID /deb-target

USER $USERNAME

ENV GOPATH=/go

RUN ls -la

RUN make

CMD ["fpm", "--version"]
