[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

1. [Usage](/README.md#usage)
2. [Installation](/docs/setup.md#installation)
3. [A note on git hooks](/README.md#a-note-on-git-hooks)
4. [TODOs](/README.md#todos)
5. [Similar Projects](/README.md#similar-projects)

## Usage
#### Setup some alias -> co-author assignments for convenience
```bash
git team assignments add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current assignments use:
```bash
git team assignments
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

## TODOs
- **internal quality:** refactor `Dockerfile`s and `Makefile` for a better development experience
- **internal quality:** get rid of command-specific functionality in `gitconfig` module

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)

