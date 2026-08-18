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

FROM bats/bats:1.13.0@sha256:6e4b9369468b7f3fd8f402ac6cc8ea7b2e4903eae28d08785f31a0245eb51a44

RUN apk add --no-cache git

COPY --from=git-team /git-team-source/target/bin/git-team /usr/local/bin/git-team

COPY src/command/enable/hookscript/prepare-commit-msg-git-team.sh /usr/local/bin/prepare-commit-msg-git-team
RUN chmod +x /usr/local/bin/prepare-commit-msg-git-team

ENV USERNAME=git-team-hookscript-test
RUN adduser -D ${USERNAME}
USER ${USERNAME}
