# Usage:
# ------
# DOCKER_BUILDKIT=true docker build -t hekmekk/git-team-zsh-completion -f zsh-completion.Dockerfile .
# docker run --rm -ti -v `pwd`/src/command/completion/script/zsh.sh:/src/git_team_completion.sh:ro hekmekk/git-team-zsh-completion
#
# % source git_team_completion.sh
# % git-team <tab> | git team <tab>

FROM alpine:3.16.2@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870

ENV RUNNING_IN_DOCKER=true
ENV USERNAME=git-team-zsh-completion

RUN apk update && \
    apk add git zsh vim zsh-autosuggestions zsh-syntax-highlighting bind-tools curl go && \
    rm -rf /var/cache/apk/*

RUN adduser -D ${USERNAME}

RUN sh -c "$(wget https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh -O -)"

RUN echo "source /usr/share/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" >> ~/.zshrc && \
    echo "source /usr/share/zsh/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" >> ~/.zshrc

ENV GOPATH=/go

RUN mkdir -p ${GOPATH} && chown -R ${USERNAME} /go && chmod -R 2750 ${GOPATH}

USER ${USERNAME}

ENV PATH=${GOPATH}/bin:$PATH

WORKDIR /git-team

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY src ./src

RUN go mod download
RUN go install ./...

RUN echo "autoload -U +X compinit && compinit" >> /home/${USERNAME}/.zshrc

VOLUME /src/git_team_completion.sh

WORKDIR /src

ENTRYPOINT /bin/zsh
