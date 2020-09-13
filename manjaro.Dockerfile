# docker build -t git-team-manjaro-debug -f manjaro.Dockerfile .
# docker run -it git-team-manjaro-debug /bin/bash

FROM manjarolinux/base:latest

LABEL maintainer Rea Sand <hekmek@posteo.de>

ARG USER_NAME=git-team-debug
ARG UID=1000
ARG GID=1000

RUN groupadd -g $GID $USER_NAME
RUN useradd -m -u $UID -g $GID -s /bin/bash $USER_NAME
RUN usermod -aG wheel $USER_NAME
RUN echo "%wheel ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN pacman --noconfirm -Sy git
RUN git clone https://aur.archlinux.org/yay-git.git /opt/yay-git
RUN chown -R $USER_NAME:$USER_NAME /opt/yay-git

RUN pacman --noconfirm -Sy make glibc gcc go gettext awk
RUN pacman --noconfirm -Sy sudo fakeroot

USER $USER_NAME

WORKDIR /opt/yay-git
RUN makepkg --noconfirm -si

RUN yay --noconfirm -Syu
RUN yay -S --answerdiff=None --nocleanmenu --nodiffmenu git-team-git

