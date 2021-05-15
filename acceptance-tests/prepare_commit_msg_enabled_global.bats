#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/prepare-commit-msg-enabled-repo-local
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
	mkdir -p $REPO_PATH
	cd $REPO_PATH
	git config --global init.defaultBranch main
	git init
	git config user.name "$USER_NAME"
	git config user.email "$USER_EMAIL"
}

teardown() {
	cd -
	rm -rf $REPO_PATH
	rm /root/.gitconfig
}

@test "when git-team is enabled (scope: global), a merge commit should have co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git checkout -b some-branch
	echo 'some-branch' >> README.md
	git commit -am 'added line to README.md'

	git checkout main

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
	# Note: fast-forward will not result in the prepare-commit-msg hook being triggered
	git merge --no-ff some-branch

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 --regexp '^Merge:'
	assert_line --index 2 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 3 --regexp '^Date:.+'
	assert_line --index 4 --regexp "\s+Merge branch 'some-branch'"
	refute_line --index 5 --regexp '\w+'
	assert_line --index 6 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 8 --regexp '\s+Co-authored-by: C <c@x.y>'

	/usr/local/bin/git-team disable
}

@test "when git-team is enabled (scope: global), a squash merge commit should have co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git checkout -b some-branch
	echo 'some-branch' >> README.md
	git commit -am 'added line to README.md'

	git checkout main

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
	git merge --squash some-branch
	git commit -m "squashed"

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+squashed'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'README.md'

	/usr/local/bin/git-team disable
}

@test "when git-team is enabled (scope: global), an amended commit message should not have any co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
	git config core.editor echo
	git commit --amend

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+initial commit'
	assert_line --index 4 --regexp 'README.md'
	refute_line --index 5 --regexp '\w+'

	/usr/local/bin/git-team disable
}

