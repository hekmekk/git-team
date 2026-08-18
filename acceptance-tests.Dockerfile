FROM golang:1.25-alpine@sha256:1e0126852075c9c60731c8ba49088448b91f63e2aed97ca9d1a9791622a05946 AS git-team

RUN apk add --no-cache make

RUN mkdir /git-team-source
WORKDIR /git-team-source

ENV GOPATH=/go

COPY go.* ./
RUN go mod download

COPY src ./src
COPY main.go .

COPY Makefile .
COPY script ./script
RUN make build

# ----------------------------------------------------------------- #

FROM bats/bats:1.13.0@sha256:6e4b9369468b7f3fd8f402ac6cc8ea7b2e4903eae28d08785f31a0245eb51a44

RUN apk add --no-cache git

COPY --from=git-team /git-team-source/target/bin/git-team /usr/local/bin/git-team

ENV USERNAME=git-team-acceptance-test
RUN adduser -D ${USERNAME}
USER ${USERNAME}
