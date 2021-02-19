FROM golang:1.14-stretch

LABEL maintainer Rea Sand <hekmek@posteo.de>

ARG USERNAME=git-team-run
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USERNAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USERNAME

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install man-db

COPY bash_completion /source/bash_completion
COPY git-hooks /source/git-hooks
COPY src /source/src
COPY cmd /source/cmd
COPY go.mod /source
COPY Makefile /source

WORKDIR /source

ENV GOPATH=/go

RUN make

RUN make install

USER $USERNAME

ENTRYPOINT ["/bin/git-team"]
