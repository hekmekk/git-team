[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for creating git commit templates provisioned with one or more co-authors.

## Installation
#### Build from Source
```bash
go get github.com/hekmekk/git-team
cd $GOPATH/github.com/hekmekk/git-team
make
sudo make install
```
#### With an [ansible playbook](../contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
```

## Usage

#### Setup some aliases
```bash
git team add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current aliases use:
```bash
git team list
```

To remove an alias use:
```bash
git team rm noujz
```

#### Provision a commit template
This alias (along with others) can then be used as an argument to the `enable` command and will be resolved while parsing the command line.
```bash
git team [enable] noujz
```
Apart from one or more aliases, you may provide a properly formatted co-author aswell.
```bash
git team [enable] noujz <alias1> ... <aliasN> "Mr. Green <green@mr.se>"
```

#### Commit some
Just use `git commit`. Please note that templates don't affect `-m`.

#### Back to being a loner
```bash
git team disable
```

## Uninstall
```bash
sudo make purge
```

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)
