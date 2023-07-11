[![Build Status](https://img.shields.io/github/actions/workflow/status/hekmekk/git-team/verify.yaml?branch=main&logo=github&style=for-the-badge)](https://github.com/hekmekk/git-team/actions)
[![Github Release Downloads](https://img.shields.io/github/downloads/hekmekk/git-team/total.svg?label=DOWNLOADS&logo=github&style=for-the-badge)](https://github.com/hekmekk/git-team/releases)

[![Legacy Code: O++ S++ I+ C>+ E+++ M+ V+ D++ Ar>A](https://img.shields.io/badge/Legacy%20Code-O%2B%2B%20S%2B%2B%20I%2B%20C%3E%2B%20E%2B%2B%2B%20M%2B%20V%2B%20D%2B%2B%20Ar%3EA-informational)](https://gitlab.com/resident-uhlig/legacy-code)

[**:heart: Donate**](/docs/donate.md#heart-donate)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

1. [Installation](/docs/installation.md)
2. [Usage](/README.md#usage)
3. [Configuration](/README.md#configuration)
4. [A note on git hooks](/README.md#a-note-on-git-hooks)
5. [Similar Projects](/README.md#similar-projects)

## Usage
### Setup some alias -> co-author assignments for convenience
```shell
git team assignments add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current assignments use:
```shell
git team assignments
```

### Set active co-authors
Apart from one or more aliases, you may provide a properly formatted co-author to the `enable` command as well.
This will activate git team globally, so that you can seemlessly switch between repositories while collaborating.
If you prefer per repository activation, you can set [the corresponding config option](/README.md#configuration).

```shell
git team enable noujz <alias1> ... <aliasN> "Mr. Green <green@mr.se>"
```

### Commit some
Just use `git commit` or `git commit -m <msg>`.

### Disable git team
```shell
git team disable
```

## Configuration
See `git team config -h` on how to configure git team.

| option             | type     | values                 | default  | description                                                    |
| ------------------ | -------- | ---------------------- | -------- | -------------------------------------------------------------- |
| `activation-scope` | `string` | `global`, `repo-local` | `global` | set to `repo-local` to use git-team on a per repository basis. |

## A note on git hooks
git-team uses a `prepare-commit-msg` hook to inject co-authors into a commit message. This hook is installed into `${HOME}/.git-team/hooks`. When you `enable` git-team, the git config option `core.hooksPath` will be set to point to that directory. Along with the `prepare-commit-msg` hook come proxies for all the other git hooks, so that other existing repo-local hooks are still being triggered.

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)

