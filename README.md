[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

## Installation
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

## Usage

#### Setup some aliases for convenience
```bash
git team add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current alias to co-author assignments use:
```bash
git team ls
```

To remove an assignment use:
```bash
git team rm noujz
```

#### Set active co-authors
Apart from one or more aliases, you may provide a properly formatted co-author aswell.
```bash
git team [enable] noujz <alias1> ... <aliasN> "Mr. Green <green@mr.se>"
```

#### Commit some
Just use `git commit` or `git commit -m <msg>`.

#### Back to being a loner
```bash
git team disable
```

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)

## TODOs
- **internal quality:** apply pattern as exemplified by `add`
- **internal quality:** refactor `Dockerfile`s and `Makefile` for a better development experience

