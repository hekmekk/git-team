FROM golang:1.14-stretch as bats

LABEL maintainer Rea Sand <hekmek@posteo.de>

RUN mkdir /bats-source
RUN git clone https://github.com/bats-core/bats-core.git --branch v1.2.0 --single-branch /bats-source
WORKDIR /bats-source
RUN ./install.sh /usr/local

WORKDIR /
RUN mkdir /bats-libs
RUN git clone https://github.com/ztombol/bats-support /bats-libs/bats-support
RUN git clone https://github.com/ztombol/bats-assert /bats-libs/bats-assert

# ----------------------------------------------------------------- #

FROM golang:1.14-stretch as git-team

RUN mkdir /git-team-source
WORKDIR /git-team-source

ENV GOPATH=/go

COPY go.* ./
RUN go mod download

COPY src ./src
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./cmd/...

# ----------------------------------------------------------------- #

FROM golang:1.14-stretch
COPY --from=bats /usr/local/bin/bats /usr/local/bin/bats
COPY --from=bats /usr/local/libexec/bats-core /usr/local/libexec/bats-core
COPY --from=bats /bats-libs /bats-libs
COPY --from=git-team /go/bin/git-team /usr/local/bin/git-team
COPY --from=git-team /go/bin/prepare-commit-msg /usr/local/bin/prepare-commit-msg-git-team

WORKDIR /

ENTRYPOINT ["bash", "/usr/local/bin/bats"]
