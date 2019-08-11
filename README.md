[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

## Installation
#### Build from Source
The latest version of git-team has been built against go version 1.12.
```bash
make
sudo make install
```
#### With an [ansible playbook](../master/contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
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
This alias (along with others) can then be used as an argument to the `enable` command and will be resolved while parsing the command line.
```bash
git team [enable] noujz
```
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

