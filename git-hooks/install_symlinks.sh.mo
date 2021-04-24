#!/bin/sh

AVAILABLE_GIT_HOOKS="applypatch-msg commit-msg fsmonitor-watchman p4-pre-submit post-applypatch post-checkout post-commit post-index-change post-merge post-receive post-rewrite post-update pre-applypatch pre-auto-gc pre-commit pre-push pre-rebase pre-receive push-to-checkout sendemail-validate update"

for hook in ${AVAILABLE_GIT_HOOKS}; do
	ln -fs proxy.sh {{ hooks_dir }}/${hook}
done
