FROM golang:1.24-alpine@sha256:2d40d4fc278dad38be0777d5e2a88a2c6dee51b0b29c97a764fc6c6a11ca893c AS git-team

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

FROM bats/bats:1.13.0@sha256:7163cd5d4f8b1e85dfd43f388c8f481e26e8b40875536dda46ebf37a0cad4eb2

RUN apk add --no-cache git

COPY --from=git-team /git-team-source/target/bin/git-team /usr/local/bin/git-team

ENV USERNAME=git-team-acceptance-test
RUN adduser -D ${USERNAME}
USER ${USERNAME}
