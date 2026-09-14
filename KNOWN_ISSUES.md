# Known Issues

## a previous commit.template setting will be destroyed by git-team

A previous `commit.template` should be stored, analogous to `team.state.previous-hooks-path`, on `enable` and restored on `disable`.

## git config might not be the best config store

It felt neat at the time to use git config as a persistent store for git-team, as it is designed as a git subcommand.

There are some drawbacks however. First, git-team aliases really pollute the global git config. Second, as described in [prepare-commit-msg-git-team.sh](src/command/enable/hookscript/prepare-commit-msg-git-team.sh), the current activation scope implementation has a conceptual problem: Changing scope outside of a git repo leaves the respective local git config unreachable.

If git-team only ever set global git options `core.hookspath` and `commit.template`, used its own storage backend (file, sqlite, ...) to store aliases and the state (global and repo-local) these issues could be mitigated while resolving the complexity to translate git-team activation-scope to git config scope.

## complected code

Over time, both prod as well as test code became quite complected.
The implementation now feels more in the way than helpful.

### go tests

The decision to enable easy mocking in tests is responsible for quite some additional abstraction.
If go tests are transformed into either true unit tests for logic only (utilizing e.g. functional core),
or use [nullables](https://www.jamesshore.com/v2/projects/nullables/testing-without-mocks), test should become a lot more readable.

### end to end tests

If coverage is lost, additional bats end to end tests are required.

### over-engineering

The command/policy/event pattern was an experiment at the time.
It's probably wise to revisit that decision to opt for a more simple approach.

Simplification of tests (see above) might already do some heavy lifting here.
