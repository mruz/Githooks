#!/bin/sh

TEST_DIR=$(dirname "$0")

cat <<EOF | docker build --force-rm -t githooks:alpine-lfs-go-whitespace-base -
FROM golang:1.15.6-alpine
RUN apk add git git-lfs --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/main --allow-untrusted
RUN apk add bash
RUN mkdir -p "/root/whitespace folder"
ENV HOME="/root/whitespace folder"
EOF

export ADDITIONAL_INSTALL_STEPS='
# add a space in paths and wrap in double-quotes
RUN find /var/lib/tests -name "*.sh" -exec sed -i -E "s|/tmp/test([0-9.]+)|\"/tmp/test \1\"|g" {} \;
# remove the double-quotes if the path is the only thing on the whole line
RUN find /var/lib/tests -name "*.sh" -exec sed -i -E "s|^\"/tmp/test([^\"]+)\"|/tmp/test\1|g" {} \;
'

exec sh "$TEST_DIR"/exec-tests-go.sh 'alpine-lfs-go-whitespace' "$@"