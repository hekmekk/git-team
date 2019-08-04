FROM golang:1.12-stretch as bats

LABEL maintainer Rea Sand <hekmek@posteo.de>

RUN mkdir /bats-source

RUN git clone https://github.com/bats-core/bats-core.git /bats-source

WORKDIR /bats-source

RUN ./install.sh /usr/local

# ----------------------------------------------------------------- #

FROM golang:1.12-stretch as git-team

RUN mkdir /git-team-source
WORKDIR /git-team-source

ENV GOPATH=/go

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/src src/main.go src/main_test.go ./
COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /usr/local/bin/git-team

# ----------------------------------------------------------------- #

FROM golang:1.12-stretch
COPY --from=bats /usr/local/bin/bats /usr/local/bin/bats
COPY --from=bats /usr/local/libexec/bats-core /usr/local/libexec/bats-core
COPY --from=git-team /usr/local/bin/git-team /usr/local/bin/git-team

WORKDIR /

ENTRYPOINT ["bash", "/usr/local/bin/bats"]
