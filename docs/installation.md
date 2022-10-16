# Installation
1. [Install git-team](#install-git-team)
    1. [Golang](#golang)
    2. [Homebrew](#homebrew)
    3. [Debian / Ubuntu](#debian--ubuntu)
    4. [RedHat / CentOS / Fedora](#redhat--centos--fedora)
    5. [Void Linux](#void-linux)
    6. [Arch Linux](#arch-linux)
    7. [Nix](#nix)
    8. [Make](#make)
2. [Setup Completion](#setup-completion)

## Install git-team
### Golang
For releases since v1.6.0.
```shell script
go install github.com/hekmekk/git-team@latest
```

### Homebrew
See [homebrew-git-team](https://github.com/hekmekk/homebrew-git-team) for the formula.
1. Add tap

```shell
brew tap hekmekk/git-team
```

2. Install git-team

Install stable release. Use `--HEAD` in case you want to install from the latest commit.
```shell
brew install git-team
```

### Debian / Ubuntu
#### apt / aptitude
0. Pre-requisites
```shell
sudo apt install curl gnupg lsb-release
```

The debian package is made available via [apt-sourc.es](https://apt-sourc.es).

1. Add *apt-sourc.es* GPG Key
```shell
curl -fsSL https://apt-sourc.es/admin/gpg.asc | sudo apt-key add -
```

2. Setup the `apt` repository
```shell
echo "deb [arch=amd64] https://apt-sourc.es/deb/hekmekk/git-team stable main" | sudo tee /etc/apt/sources.list.d/git-team.list
```

3. Update the `apt` package index and install `git-team`
```shell
sudo apt update && sudo apt install git-team
```

#### [ansible playbook](../master/contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
```

#### manually
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Install it manually
```shell
sudo dpkg -i /path/to/downloaded/release.deb
```

### RedHat / CentOS / Fedora
#### [ansible playbook](../master/contrib/ansible/roles/git-team/tasks/main.yml)
```
ansible-playbook git-team.yml --ask-become-pass
```

#### manually
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Import git-team signing key
```shell
curl --silent https://api.github.com/users/hekmekk/gpg_keys | jq -r '.[] | select(.key_id == "12BB70967049E845") | .raw_key' > git-team-signing-key.asc
```

```shell
sudo rpm --import git-team-signing-key.asc
```

3. Check if Signature is ok
```shell
rpm -v --checksig /path/to/downloaded/release.rpm
```

4. Install
```shell
sudo rpm -i /path/to/downloaded/release.rpm
```

### Void Linux
[@steinex](https://github.com/steinex) maintains the [template](https://github.com/void-linux/void-packages/blob/master/srcpkgs/git-team/template).
```shell
sudo xbps-install git-team
```

### Arch Linux
[@lockejan](https://github.com/lockejan) maintains the [AUR](https://aur.archlinux.org/packages/git-team-git/).
```shell
yay git-team-git
```

### [Nix](https://nixos.org)

```nix
{ pkgs, ... }:
{
  # Either for all users
  environment.systemPackages = with pkgs; [ git-team ];

  # Or for an explicit user
  users.users."youruser".packages = with pkgs; [ git-team ];
}
```

### Make
The latest version of git-team has been built against go version 1.18.
```shell
make
sudo make install
```

## Setup Completion
Check if a script is available for your shell
```shell
git-team completion
```

Source the appropriate script, e.g.:
```bash
source <(git-team completion bash)
```

