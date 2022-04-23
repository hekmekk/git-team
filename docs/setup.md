# Installation
## OSX
### via Homebrew
See [homebrew-git-team](https://github.com/hekmekk/homebrew-git-team) for the formula.
1. Add tap

```bash
brew tap hekmekk/git-team
```

2. Install git-team

Install stable release. Use `--HEAD` in case you want to install from the latest commit.
```bash
brew install git-team
```

## Debian / Ubuntu
### via apt / aptitude
0. Pre-requisites
```bash
sudo apt install curl gnupg lsb-release
```

#### using [apt-sourc.es](https://apt-sourc.es)
1. Add *apt-sourc.es* GPG Key
```bash
curl -fsSL https://apt-sourc.es/admin/gpg.asc | sudo apt-key add -
```

2. Setup the `apt` repository
```bash
echo "deb [arch=amd64] https://apt-sourc.es/deb/hekmekk/git-team stable main" | sudo tee /etc/apt/sources.list.d/git-team.list
```

3. Update the `apt` package index and install `git-team`
```bash
sudo apt update && sudo apt install git-team
```

### via an [ansible playbook](../master/contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
```

### manually
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Install it manually
```bash
sudo dpkg -i /path/to/downloaded/release.deb
```

## RedHat / CentOS / Fedora
### manually
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Install it manually
```bash
sudo rpm -i /path/to/downloaded/release.rpm
```

## Void Linux
[@steinex](https://github.com/steinex) maintains the [template](https://github.com/void-linux/void-packages/blob/master/srcpkgs/git-team/template).
```bash
sudo xbps-install git-team
```

## Arch Linux
[@lockejan](https://github.com/lockejan) maintains the [AUR](https://aur.archlinux.org/packages/git-team-git/).
```bash
yay git-team-git
```

## Install declaratively via [Nix](https://nixos.org)

```nix
{ pkgs, ... }:
{
  # Either for all users
  environment.systemPackages = with pkgs; [ git-team ];

  # Or for an explicit user
  users.users."youruser".packages = with pkgs; [ git-team ];
}
```

## Build from source
The latest version of git-team has been built against go version 1.18.
```bash
make
sudo make install
```

## Install from source
For releases since v1.6.0.
```shell script
go install github.com/hekmekk/git-team@latest
source <(git-team completion bash)
```

