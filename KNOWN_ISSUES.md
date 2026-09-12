# Known Issues

## A previous commit.template setting will be destroyed by git-team

A previous `commit.template` should be stored, analogous to `team.state.previous-hooks-path`, on `enable` and restored on `disable`.

## Complected design with too many overlapping tests

Over time, both prod as well as test code became quite complected.
The implementation now feels more in the way than helpful.

The decision to enable easy mocking in tests is responsible for quite some additional abstraction.
If go tests are transformed into true unit tests, mocking might become obsolete.
This should probably be addressed first and might require some more modularization.
If coverage is lost, additional bats end to end tests are required.

The command/policy/event pattern was an experiment at the time.
It's probably wise to revisit that decision to opt for a more simple approach.
