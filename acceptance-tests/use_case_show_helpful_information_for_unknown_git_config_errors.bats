#!/usr/bin/env bats

REPO_PATH=/tmp/repo/non-writable-gitconfig
USER_NAME=git-team-acceptance-test
USER_EMAIL=acc@git.team

setup() {
        bats_load_library bats-support
        bats_load_library bats-assert

        git config --global init.defaultBranch main

        # lookup git config paths: /usr/bin/env git config --list --show-origin --show-scope
        chmod -x /home/git-team-acceptance-test/

        mkdir -p $REPO_PATH
        cd $REPO_PATH
        touch THE_FILE

        git init
        git config user.name "$USER_NAME"
        git config user.email "$USER_EMAIL"
}

teardown() {
        cd -
        rm -rf $REPO_PATH

        chmod +x /home/git-team-acceptance-test/
        rm /home/git-team-acceptance-test/.gitconfig
}

@test "use case: the user should see helpful debug information in case of an unknown gitconfig error (e.g. non-writable config)" {
        # Note: scope is irrelevant here, because adding assignments is always done with scope=global
        run /usr/local/bin/git-team assignments add noujz "Mr. Noujz <noujz@mr.se>"

        assert_failure
        assert_line --index 0 'error: failed to add alias: unknown gitconfig error: exit status 255'
        assert_line --index 1 'git-config command:'
        assert_line --index 2 '/usr/bin/env git config --global --replace-all team.alias.noujz Mr. Noujz <noujz@mr.se>'
        assert_line --index 3 'git-config output:'
        assert_line --index 4 "warning: unable to access '/home/git-team-acceptance-test/.gitconfig': Permission denied"
        assert_line --index 5 "warning: unable to access '/home/git-team-acceptance-test/.config/git/config': Permission denied"
        assert_line --index 6 "error: could not lock config file /home/git-team-acceptance-test/.gitconfig: Permission denied"
}

