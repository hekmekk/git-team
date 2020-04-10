[![Build Status](https://travis-ci.org/hekmekk/git-team.svg?branch=master)](https://travis-ci.org/hekmekk/git-team)

# git-team

Command line interface for managing and enhancing `git commit` messages with co-authors.

1. [Installation](/docs/setup.md#installation)
2. [Usage](/README.md#usage)
3. [Configuration](/README.md#configuration)
4. [A note on git hooks](/README.md#a-note-on-git-hooks)
5. [TODOs](/README.md#todos)
6. [Similar Projects](/README.md#similar-projects)

## Usage
### Setup some alias -> co-author assignments for convenience
```bash
git team assignments add noujz "Mr. Noujz <noujz@mr.se>"
```

To review your current assignments use:
```bash
git team assignments
```

### Set active co-authors
Apart from one or more aliases, you may provide a properly formatted co-author to the `enable` command as well.
This will activate git team globally, so that you can seemlessly switch between repositories while collaborating.
If you prefer per repository activation, you can set [the corresponding config option](/README.md#configuration).

```bash
git team [enable] noujz <alias1> ... <aliasN> "Mr. Green <green@mr.se>"
```

### Commit some
Just use `git commit` or `git commit -m <msg>`.

### Disable git team
```bash
git team disable
```

## Configuration
See `git team config -h` on how to configure git team.

| option           | values             | default | affected commands   |
| ---------------- | ------------------ | ------- | ------------------- |
| activation-scope | global, repo-local | global  | `enable`, `disable` |

## A note on git hooks
git-team uses a `prepare-commit-msg` hook to inject co-authors into a commit message. This hook is installed into `/usr/local/etc/git-team/hooks`. When you `enable` git-team, the git config option `core.hooksPath` will be set to point to that directory. Along with the `prepare-commit-msg` hook come proxies for all the other git hooks, so that other existing repo-local hooks are still being triggered.

## TODOs
- **internal quality:** refactor `src/command/adapter/alias_resolver` to `assignments_repository` residing in `core`
- **internal quality:** consequently separate domain-specific from technological functionality
- **internal quality:** consequently use domain-specific first and technology second packaging
- **internal quality:** define public interfaces where they're missing
- **internal quality:** introduce public interfaces package(s)
- **internal quality:** refactor `Dockerfile`s and `Makefile` for a better development experience

## Similar projects
- [git mob](https://www.npmjs.com/package/git-mob)

