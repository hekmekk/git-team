FROM golang:1.24-alpine AS git-team

RUN mkdir /git-team-source
WORKDIR /git-team-source

ENV GOPATH=/go

COPY go.* ./
RUN go mod download

COPY src ./src
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go install ./...

# ----------------------------------------------------------------- #

FROM bats/bats:1.13.0

RUN apk add --no-cache git

COPY --from=git-team /go/bin/git-team /usr/local/bin/git-team

ENV USERNAME=git-team-acceptance-test
RUN adduser -D ${USERNAME}
USER ${USERNAME}
