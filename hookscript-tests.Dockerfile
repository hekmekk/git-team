FROM alpine/git:v2.45.2 AS bats

LABEL maintainer="Rea Sand <hekmek@posteo.de>"

RUN apk add --no-cache bash

RUN mkdir /bats-source
RUN git clone https://github.com/bats-core/bats-core.git --branch v1.2.0 --single-branch /bats-source
WORKDIR /bats-source
RUN ./install.sh /usr/local

WORKDIR /
RUN mkdir /bats-libs
RUN git clone https://github.com/ztombol/bats-support /bats-libs/bats-support
RUN git clone https://github.com/ztombol/bats-assert /bats-libs/bats-assert

# ----------------------------------------------------------------- #

FROM golang:1.23-alpine AS git-team

RUN mkdir /git-team-source
WORKDIR /git-team-source

ENV GOPATH=/go

COPY go.* ./
RUN go mod download

COPY src ./src
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go install ./...

# ----------------------------------------------------------------- #

FROM golang:1.23-alpine

RUN apk add --no-cache bash git ncurses

COPY --from=bats /usr/local/bin/bats /usr/local/bin/bats
COPY --from=bats /usr/local/libexec/bats-core /usr/local/libexec/bats-core
COPY --from=bats /bats-libs /bats-libs
COPY --from=git-team /go/bin/git-team /usr/local/bin/git-team
COPY src/command/enable/hookscript/prepare-commit-msg-git-team.sh /usr/local/bin/prepare-commit-msg-git-team

RUN chmod +x /usr/local/bin/prepare-commit-msg-git-team

WORKDIR /

ENTRYPOINT ["bash", "/usr/local/bin/bats"]
