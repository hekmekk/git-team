# Installation
## OSX
### via [Homebrew](https://brew.sh)
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

#### using [bintray](https://bintray.com)
1. Add bintray GPG Key
```bash
curl -fsSL https://api.bintray.com/users/hekmekk/keys/gpg/public.key | sudo apt-key add -
```

2. Setup the `apt` repository
```bash
echo "deb [arch=amd64] https://dl.bintray.com/hekmekk/git-team $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/git-team.list
```

3. Update the `apt` package index and install `git-team`
```bash
sudo apt update && sudo apt install git-team
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
### via yum
#### using [bintray](https://bintray.com)
1. Add bintray GPG Key
```bash
rpm --import https://api.bintray.com/users/hekmekk/keys/gpg/public.key
```

2. Setup the `yum` repository
```bash
echo "[git-team]
name=git-team
enabled=1
baseurl=https://dl.bintray.com/hekmekk/rpm
gpgcheck=1
gpgkey=https://api.bintray.com/users/hekmekk/keys/gpg/public.key" | sudo tee /etc/yum.repos.d/git-team.repo
```

3. Install `git-team`
```bash
sudo yum install git-team
```

### manually
1. Download the [latest release](https://github.com/hekmekk/git-team/releases/latest)

2. Install it manually
```bash
sudo rpm -i /path/to/downloaded/release.rpm
```

## Void Linux
```bash
sudo xbps-install git-team
```

## Arch Linux
Install from the [AUR](https://aur.archlinux.org/packages/git-team-git/)
```bash
yay git-team-git
```

## Build from source
The latest version of git-team has been built against go version 1.14.
```bash
make
sudo make install
```

