# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [1.5.2] - 2020-11-10
### Fixed
- Parameters are being passed to existing git hooks.

## [1.5.1] - 2020-10-30
### Fixed
- Usage section for different commands and subcommands.

## [1.5.0] - 2020-10-30
### Added
- New flag `--all|-A` for the `enable` command to include all known co-authors.

### Changed
- Shell completion is now done within the application itself.

### Deprecated
- `enable` as a default command

## [1.4.1] - 2020-09-13
### Fixed
- Don't append co-authors to a commit message if they are part of it already. This may happen when co-authors have been added manually or when both the commit template as well as the `prepare-commit-msg` hook take effect (e.g.: IDEs reading from the commit template and writing the entire content (including co-authors) back as a single commit message).

## [1.4.0] - 2020-05-10
### Added
- New command `config` has been introduced to view and edit the configuration
- The `activation-scope` can now be configured (options: `global` (default), `repo-local`) via `git-team config`

## [1.3.8] - 2019-10-26
### Added
- Co-authors are appended for `git merge [--squash]`

## [1.3.7] - 2019-10-04
### Changed
- Minor adjustments to the `bash_comletion` script

## [1.3.6] - 2019-10-02
### Added
- New flag `-f|--force-override` for `assignments add` subcommand

## [1.3.5] - 2019-09-19
### Added
- Command `assignments` with subcommands (`add`, `rm` and `ls`)

### Deprecated
- Commands `add`, `rm` and `ls`

### Changed
- Command line output format
- Configuration is now stored within `git config`

### Fixed
- Signal error when trying to remove a non-existing assignment

## [1.3.4] - 2019-09-04
### Changed
- Proxy scripts for repo-local git hooks are now symlinked

## [1.3.3] - 2019-09-03
### Fixed
- Install proxy scripts for repo-local git hooks so that they won't be disabled when git-team is enabled

## [1.3.2] - 2019-08-30
### Fixed
- Show an empty list when there are no assignments instead of failing with an error

## [1.3.1] - 2019-08-20
### Changed
- Refactor Makefile
- Add section on homebrew as an installation option
- Install git-team to `/usr/bin`
- build packages for multiple targets (deb and rpm)

## [1.3.0] - 2019-08-11
### Added
- Support `git commit -m <msg>` via a `prepare-commit-msg` hook
- Ask the user if an existing assignment should be overridden

### Changed
- The `list` command is now an alias for `ls` and will no longer be suggested via auto-completion

## [1.2.1] - 2019-08-07
### Fixed
- Make the `disable` command work reliably. It would occasionally not have any effect at all.

## [1.2.0] - 2019-07-30
### Changed
- Always sort lists presented to the user (aliases, co-authors)
- Ignore duplicates when running `enable` command

## [1.1.3] - 2019-07-28
### Changed
- Implementation of `add` command

## [1.1.2] - 2019-07-26
### Fixed
- Persistence of current status (enabled, disabled)

## [1.1.1] - 2019-07-26
### Changed
- Adjust code structure to increase test coverage

## [1.1.0] - 2019-07-14
### Added
- `ls` as an alias for `list` command

## [1.0.1] - 2019-06-30
### Fixed
- Show an empty list when there are no assignments instead of failing with an error

## [1.0.0] - 2019-04-17
### Changed
- Make this go v1.12 compliant

## [0.3.0] - 2018-10-28
### Added
- The ability to manage co-authors by assigning aliases to them

## [0.2.0] - 2018-10-24
### Fixed
- Installation via `make install` on osx

## [0.1.0] - 2018-08-06
### Changed
- Remove dependency on `git2go` and therefore `libgit2`

## [0.0.1] - 2018-07-17
### Added
- Append co-authors to a `git commit` message by means of a commit message template file

[Unreleased]: https://github.com/hekmekk/git-team/compare/v1.5.2...HEAD
[1.5.2]: https://github.com/hekmekk/git-team/compare/v1.5.1...v1.5.2
[1.5.1]: https://github.com/hekmekk/git-team/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/hekmekk/git-team/compare/v1.4.1...v1.5.0
[1.4.1]: https://github.com/hekmekk/git-team/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/hekmekk/git-team/compare/v1.3.8...v1.4.0
[1.3.8]: https://github.com/hekmekk/git-team/compare/v1.3.7...v1.3.8
[1.3.7]: https://github.com/hekmekk/git-team/compare/v1.3.6...v1.3.7
[1.3.6]: https://github.com/hekmekk/git-team/compare/v1.3.5...v1.3.6
[1.3.5]: https://github.com/hekmekk/git-team/compare/v1.3.4...v1.3.5
[1.3.4]: https://github.com/hekmekk/git-team/compare/v1.3.3...v1.3.4
[1.3.3]: https://github.com/hekmekk/git-team/compare/v1.3.2...v1.3.3
[1.3.2]: https://github.com/hekmekk/git-team/compare/v1.3.1...v1.3.2
[1.3.1]: https://github.com/hekmekk/git-team/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/hekmekk/git-team/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/hekmekk/git-team/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/hekmekk/git-team/compare/v1.1.3...v1.2.0
[1.1.3]: https://github.com/hekmekk/git-team/compare/v1.1.2...v1.1.3
[1.1.2]: https://github.com/hekmekk/git-team/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/hekmekk/git-team/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/hekmekk/git-team/compare/v0.0.1...v1.1.0
[1.0.1]: https://github.com/hekmekk/git-team/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/hekmekk/git-team/compare/v0.3.0...v1.0.0
[0.3.0]: https://github.com/hekmekk/git-team/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/hekmekk/git-team/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/hekmekk/git-team/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/hekmekk/git-team/releases/tag/v0.0.1

