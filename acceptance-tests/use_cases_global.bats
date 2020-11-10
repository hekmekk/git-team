#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/use-cases-global

setup() {
	cp /usr/local/bin/prepare-commit-msg /usr/local/etc/git-team/hooks/prepare-commit-msg

	mkdir -p $REPO_PATH
	cd $REPO_PATH
	touch THE_FILE

	git init
	git config user.name git-team-acceptance-test
	git config user.email foo@bar.baz

	/usr/local/bin/git-team config activation-scope global
}

teardown() {
	/usr/local/bin/git-team disable

	cd -
	rm -rf $REPO_PATH
}

@test "use case: (scope: global) an existing repo-local git hook should be respected" {
	echo -e '#!/bin/sh\necho "commit-msg hook triggered with params: $@"\nexit 1' > $REPO_PATH/.git/hooks/commit-msg
	chmod +x $REPO_PATH/.git/hooks/commit-msg

	/usr/local/bin/git-team enable 'A <a@x.y>'

	git add -A
	run git commit -m "test"

	assert_failure
	assert_line --index 0 'commit-msg hook triggered with params: .git/COMMIT_EDITMSG'
}

@test "use case: (scope: global) when git-team is enabled then 'git commit -m' should have the respective co-authors injected" {
	/usr/local/bin/git-team enable 'B <b@x.y>' 'A <a@x.y>' 'C <c@x.y>'

	git add -A
	git commit -m "test"

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <foo@bar.baz>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'THE_FILE'
}

@test "use case: (scope: global) when git-team is enabled then 'git commit -m' should not result in interference with existing co-authors" {
	/usr/local/bin/git-team enable 'B <b@x.y>' 'A <a@x.y>' 'C <c@x.y>'

	git add -A
	git commit -F- <<EOF
test

Co-authored-by: D <d@x.y>
EOF

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <foo@bar.baz>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: D <d@x.y>'
	assert_line --index 6 'THE_FILE'
}

@test "use case: (scope: global) when git-team is disabled then 'git commit -m' should not have any co-authors injected" {
	/usr/local/bin/git-team disable

	git add -A
	git commit -m "test"

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <foo@bar.baz>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	assert_line --index 4 'THE_FILE'
	refute_output --partial 'Co-authored-by:'
}
