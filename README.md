[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

## Usage

#### Setup some aliases for convenience
```bash
git team add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current alias to co-author assignments use:
```bash
git team ls
```

#### Set active co-authors
Apart from one or more aliases, you may provide a properly formatted co-author aswell.
```bash
git team [enable] noujz <alias1> ... <aliasN> "Mr. Green <green@mr.se>"
```

#### Commit some
Just use `git commit` or `git commit -m <msg>`.

#### Disable git team
```bash
git team disable
```

## A note on git hooks
git-team uses a `prepare-commit-msg` hook to inject co-authors into a commit message. This hook is installed into `/usr/local/etc/git-team/hooks`. When you `enable` git-team, the git config option `core.hooksPath` will be set to point to that directory. Along with the `prepare-commit-msg` hook come proxies for all the other git hooks, so that other existing repo-local hooks are still being triggered.

## Installation
#### via [Homebrew](https://brew.sh)
1. Add tap

```bash
brew tap hekmekk/git-team
```

2. Install git-team

Install stable release. Use `--HEAD` in case you want to install from the latest commit.
```bash
brew install git-team
```

#### via [apt-sourc.es](https://apt-sourc.es)
1. Add *apt-sourc.es* GPG Key
```bash
curl https://apt-sourc.es/admin/gpg.asc | sudo apt-key add -
```

2. Setup the `apt` repository
```bash
echo "deb [arch=amd64] https://apt-sourc.es/deb/hekmekk/git-team stable main" | sudo tee /etc/apt/sources.list.d/git-team.list
```

3. Update the `apt` package index and install `git-team`
```bash
sudo apt update && sudo apt install git-team
```

#### via an [ansible playbook](../master/contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
```

#### Download a Release
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Install it manually
```bash
sudo dpkg -i /path/to/downloaded/release.deb
```

#### Build from Source
The latest version of git-team has been built against go version 1.12.
```bash
make
sudo make install
```

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)

## TODOs
- **fix**: if `core.hooksPath` is set already, symlink git-team `prepare-commit-msg` there (fail if it exists already)
- **internal quality:** apply pattern as exemplified by `add`
- **internal quality:** refactor `Dockerfile`s and `Makefile` for a better development experience

