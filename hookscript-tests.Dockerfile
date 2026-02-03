FROM golang:1.24-alpine AS git-team

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

FROM bats/bats:1.13.0

RUN apk add --no-cache git

COPY --from=git-team /git-team-source/target/bin/git-team /usr/local/bin/git-team

COPY src/command/enable/hookscript/prepare-commit-msg-git-team.sh /usr/local/bin/prepare-commit-msg-git-team
RUN chmod +x /usr/local/bin/prepare-commit-msg-git-team

ENV USERNAME=git-team-hookscript-test
RUN adduser -D ${USERNAME}
USER ${USERNAME}
