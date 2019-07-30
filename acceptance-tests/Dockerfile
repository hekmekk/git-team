FROM git-team-run:latest

LABEL maintainer Rea Sand <hekmek@posteo.de>

USER root

RUN mkdir /bats-source

RUN git clone https://github.com/bats-core/bats-core.git /bats-source

WORKDIR /bats-source

RUN ./install.sh /usr/local

RUN ls -la /usr/local/bin/bats

WORKDIR /

ENTRYPOINT ["bash", "/usr/local/bin/bats"]

