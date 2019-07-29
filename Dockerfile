FROM golang:1.12-stretch

LABEL maintainer Rea Sand <hekmek@posteo.de>

ARG USERNAME=git-team-run
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USERNAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USERNAME

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install man-db

COPY bash_completion /source/bash_completion
COPY src /source/src
COPY go.mod /source
COPY main.go /source
COPY main_test.go /source
COPY Makefile /source

WORKDIR /source

ENV GOPATH=/go

RUN make

RUN make install

USER $USERNAME

ENTRYPOINT ["/usr/local/bin/git-team"]
