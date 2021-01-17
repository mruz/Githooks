#!/bin/sh
# Test:
#   Run an install without an existing template directory and refusing to set a new one up

rm -rf /usr/share/git-core/templates/hooks

echo 'n
' | "$GH_TEST_BIN/installer" --stdin

# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
    echo "! Expected to fail"
    exit 1
fi
