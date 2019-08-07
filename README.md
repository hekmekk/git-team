[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for creating git commit templates provisioned with one or more co-authors.

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

#### Build from Source
The latest version of git-team has been built against go version 1.12.
```bash
make
sudo make install
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

## TODOs
- refactor: add `assign` and `unassign` but keep `add` and `rm` as aliases for backwards compatibility
- refactor: make `ls` the command and `list` the alias for backwards compatibility
- refactor: consolidate "persistence backends" `git config` and git-team status file
- feat: make it possible to rm multiple aliases
- feat: when adding an existing alias, ask for override
- feat: use current COMMIT TEMPLATE if one exists

