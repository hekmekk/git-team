FROM golang:1.12-stretch

LABEL maintainer Rea Sand <hekmek@posteo.de>

ENV GOPATH=/go

RUN mkdir /bats-source

RUN git clone https://github.com/bats-core/bats-core.git /bats-source

WORKDIR /bats-source

RUN ./install.sh /usr/local

WORKDIR /
RUN mkdir /source
COPY src/go.mod /source
COPY src/src /source/src
COPY src/main.go /source
COPY src/main_test.go /source

WORKDIR /source

RUN mkdir /target

RUN go build -o /target/git-team

RUN install /target/git-team /usr/local/bin/git-team

WORKDIR /

ENTRYPOINT ["bash", "/usr/local/bin/bats"]
