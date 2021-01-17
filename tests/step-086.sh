#!/bin/sh
# Test:
#   Cli tool: list Githooks configuration

if ! "$GH_TEST_BIN/installer"; then
    echo "! Failed to execute the install script"
    exit 1
fi

mkdir -p "$GH_TEST_TMP/test086" && cd "$GH_TEST_TMP/test086" || exit 3

! "$GITHOOKS_INSTALL_BIN_DIR/cli" config list --local || exit 4 # not a Git repo

git init || exit 5

"$GITHOOKS_INSTALL_BIN_DIR/cli" config update --enable || exit 7
"$GITHOOKS_INSTALL_BIN_DIR/cli" config list
"$GITHOOKS_INSTALL_BIN_DIR/cli" config list | grep -q -i 'githooks.autoUpdateEnabled' || exit 8
"$GITHOOKS_INSTALL_BIN_DIR/cli" config list --global | grep -q -i 'githooks.autoUpdateEnabled' || exit 9
! "$GITHOOKS_INSTALL_BIN_DIR/cli" config list --local | grep -q -i 'githooks.autoUpdateEnabled' || exit 10
