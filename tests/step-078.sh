#!/bin/sh
# Test:
#   Cli tool: enable/disable auto updates

TEST_DIR=$(cd "$(dirname "$0")" && pwd)
# shellcheck disable=SC1090
. "$TEST_DIR/general.sh"

acceptAllTrustPrompts || exit 1

if ! "$GH_TEST_BIN/cli" installer; then
    echo "! Failed to execute the install script"
    exit 1
fi

git config --global --unset githooks.autoUpdateEnabled &&
    "$GITHOOKS_INSTALL_BIN_DIR/cli" update --enable &&
    [ "$(git config --get githooks.autoUpdateEnabled)" = "true" ] ||
    exit 1

git config --global --unset githooks.autoUpdateEnabled &&
    "$GITHOOKS_INSTALL_BIN_DIR/cli" update --disable &&
    [ "$(git config --get githooks.autoUpdateEnabled)" = "false" ] ||
    exit 1
