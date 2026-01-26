#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/use-cases-global
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
	/usr/local/bin/git-team disable

	git config --global init.defaultBranch main

	mkdir -p $REPO_PATH
	cd $REPO_PATH

	git init
	git config user.name "$USER_NAME"
	git config user.email "$USER_EMAIL"
}

teardown() {
	/usr/local/bin/git-team disable

	cd -
	rm -rf $REPO_PATH

	rm /home/git-team-acceptance-test/.gitconfig
}

@test "use case: when git-team is disabled, a merge commit should not have any co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git checkout -b some-branch
	echo 'some-branch' >> README.md
	git commit -am 'added line to README.md'

	git checkout main

	# Note: fast-forward will not result in the prepare-commit-msg hook being triggered
	git merge --no-ff some-branch

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 --regexp '^Merge:'
	assert_line --index 2 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 3 --regexp '^Date:.+'
	assert_line --index 4 --regexp "\s+Merge branch 'some-branch'"
	refute_output --partial 'Co-authored-by:'
}

@test "use case: when git-team is disabled, a squash merge commit should not have any co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git checkout -b some-branch
	echo 'some-branch' >> README.md
	git commit -am 'added line to README.md'

	git checkout main

	git merge --squash some-branch
	git commit -m "squashed"

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+squashed'
	assert_line --index 4 'README.md'
	refute_output --partial 'Co-authored-by:'
}

@test "use case: when git-team is disabled then 'git commit -m' should not have any co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md

	git commit -m "test"

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	assert_line --index 4 'README.md'
	refute_output --partial 'Co-authored-by:'
}

# TODO: test anpassen
@test "use case: when git-team is disabled, an amended commit message should not have any co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git commit --amend -m "amended"

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+amended'
	assert_line --index 4 'README.md'
	refute_output --partial 'Co-authored-by:'
}

