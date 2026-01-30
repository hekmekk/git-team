#!/usr/bin/env bats

REPO_PATH=/tmp/repo/use-cases-global
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
	bats_load_library bats-support
	bats_load_library bats-assert

	git config --global init.defaultBranch main

	mkdir -p $REPO_PATH
	cd $REPO_PATH
	touch THE_FILE

	git init
	git config user.name "$USER_NAME"
	git config user.email "$USER_EMAIL"

	/usr/local/bin/git-team config activation-scope global
}

teardown() {
	/usr/local/bin/git-team disable

	cd -
	rm -rf $REPO_PATH

	rm /home/git-team-acceptance-test/.gitconfig
}

@test "use case: (scope: global) an existing repo-local git hook should be respected - commit-msg" {
	echo -e '#!/bin/sh\necho "commit-msg hook triggered with params: $@"\nexit 1' > $REPO_PATH/.git/hooks/commit-msg
	chmod +x $REPO_PATH/.git/hooks/commit-msg

	/usr/local/bin/git-team enable 'A <a@x.y>'

	git add -A
	run git commit -m "test"

	assert_failure
	assert_line --index 0 'commit-msg hook triggered with params: .git/COMMIT_EDITMSG'
}

@test "use case: (scope: global) an existing global git hook should be respected instead of a repo-local one - commit-msg" {
	echo -e '#!/bin/sh\necho "repo-local commit-msg hook triggered with params: $@"\nexit 1' > $REPO_PATH/.git/hooks/commit-msg
	chmod +x $REPO_PATH/.git/hooks/commit-msg

	mkdir -p /tmp/non-git-team-hooks/
	echo -e '#!/bin/sh\necho "global commit-msg hook triggered with params: $@"\nexit 0' > /tmp/non-git-team-hooks/commit-msg
	chmod +x /tmp/non-git-team-hooks/commit-msg

	git config --global core.hooksPath "/tmp/non-git-team-hooks"

	/usr/local/bin/git-team enable 'A <a@x.y>'

	git add -A
	run git commit -m "test"

	assert_success
	assert_line --index 0 'global commit-msg hook triggered with params: .git/COMMIT_EDITMSG'

	git config --global --unset core.hooksPath | true
	rm -rf /tmp/non-git-team-hooks
}

@test "use case: (scope: global) an existing repo-local git hook should be respected - prepare-commit-msg" {
	echo -e '#!/bin/sh\necho "prepare-commit-msg hook triggered with params: $@"\nexit 1' > $REPO_PATH/.git/hooks/prepare-commit-msg
	chmod +x $REPO_PATH/.git/hooks/prepare-commit-msg

	/usr/local/bin/git-team enable 'A <a@x.y>'

	git add -A
	run git commit -m "test"

	assert_failure
	assert_line --index 0 'prepare-commit-msg hook triggered with params: .git/COMMIT_EDITMSG message'
}

@test "use case: (scope: global) an existing global git hook should be respected instead of a repo-local one - prepare-commit-msg" {
	echo -e '#!/bin/sh\necho "repo-local prepare-commit-msg hook triggered with params: $@"\nexit 1' > $REPO_PATH/.git/hooks/prepare-commit-msg
	chmod +x $REPO_PATH/.git/hooks/prepare-commit-msg

	mkdir -p /tmp/non-git-team-hooks/
	echo -e '#!/bin/sh\necho "global prepare-commit-msg hook triggered with params: $@"\nexit 0' > /tmp/non-git-team-hooks/prepare-commit-msg
	chmod +x /tmp/non-git-team-hooks/prepare-commit-msg

	git config --global core.hooksPath "/tmp/non-git-team-hooks"

	/usr/local/bin/git-team enable 'A <a@x.y>'

	git add -A
	run git commit -m "test"

	assert_success
	assert_line --index 0 'global prepare-commit-msg hook triggered with params: .git/COMMIT_EDITMSG message'

	git config --global --unset core.hooksPath | true
	rm -rf /tmp/non-git-team-hooks
}

@test "use case: (scope: global) when git-team is enabled then 'git commit -m' should have the respective co-authors injected" {
	/usr/local/bin/git-team enable 'B <b@x.y>' 'A <a@x.y>' 'C <c@x.y>'

	git add -A
	git commit -m "test"

	run git show --name-only
	assert_success
	assert_line --index 0 --regexp '^commit\s\w+'
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
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
	assert_line --index 1 "Author: $USER_NAME <$USER_EMAIL>"
	assert_line --index 2 --regexp '^Date:.+'
	assert_line --index 3 --regexp '\s+test'
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: D <d@x.y>'
	assert_line --index 6 'THE_FILE'
}

@test "use case: (scope: global) when git-team is enabled, a merge commit should have co-authors injected" {
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

@test "use case: (scope: global) when git-team is enabled, a squash merge commit should have co-authors injected" {
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

@test "use case: (scope: global) when git-team is enabled, an amended commit message should have co-authors injected" {
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
	refute_line --index 4 --regexp '\w+'
	assert_line --index 5 --regexp '\s+Co-authored-by: A <a@x.y>'
	assert_line --index 6 --regexp '\s+Co-authored-by: B <b@x.y>'
	assert_line --index 7 --regexp '\s+Co-authored-by: C <c@x.y>'
	assert_line --index 8 'README.md'

	/usr/local/bin/git-team disable
}
