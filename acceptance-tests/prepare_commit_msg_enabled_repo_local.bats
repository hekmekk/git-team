#!/usr/bin/env bats

load '/bats-libs/bats-support/load.bash'
load '/bats-libs/bats-assert/load.bash'

REPO_PATH=/tmp/repo/prepare-commit-msg-enabled-repo-local
REPO_CHECKSUM=$(echo -n $USER:$REPO_PATH | md5sum | awk '{ print $1 }')

setup() {
	mkdir -p $REPO_PATH
	cd $REPO_PATH
	git config --global init.defaultBranch main
	git init
	git config user.name 'git-team-acceptance-test'
	git config user.email 'acc@git.team'
	/usr/local/bin/git-team config activation-scope repo-local
	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
}

teardown() {
	/usr/local/bin/git-team disable
	/usr/local/bin/git-team config activation-scope global
	cd -
	rm -rf $REPO_PATH
	rm /root/.gitconfig
}

@test "when git-team is enabled (scope: repo-local), a merge commit should have co-authors injected" {
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	git checkout -b some-branch
	echo 'some-branch' >> README.md
	git commit -am 'added line to README.md'

	git checkout main

	git merge some-branch

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <acc@git.team>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+added line to README.md'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'README.md'

	/usr/local/bin/git-team disable
}

@test "when git-team is enabled (scope: repo-local), a squash merge commit should have co-authors injected" {
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
	assert_line --index 1 'Author: git-team-acceptance-test <acc@git.team>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+squashed'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'README.md'

	/usr/local/bin/git-team disable
}

@test "when git-team is enabled (scope: repo-local), an amended commit message should not have any co-authors injected" {
	/usr/local/bin/git-team disable
	echo '# some-repository' > README.md
	git add README.md
	git commit -m 'initial commit'

	/usr/local/bin/git-team enable 'A <a@x.y>' 'B <b@x.y>' 'C <c@x.y>'
	git config core.editor echo
	git commit --amend

	run git show --name-only

	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 'Author: git-team-acceptance-test <acc@git.team>'
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+initial commit'
	assert_line --index 4 --regexp 'README.md'
	refute_line --index 5 --regexp '\w+'

	/usr/local/bin/git-team disable
}

