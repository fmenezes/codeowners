#!/usr/bin/env bash

if [ -z "$TAG" ]
then
      echo "\$TAG is empty"
      exit 1
fi

export PREVIOUS_TAG=$(git tag --sort=refname | grep $TAG -B 1 | grep -m1 "")

cat <<EOF
## Changes
$(git log $TAG...$PREVIOUS_TAG --format="- [%h](../../commit/%h) %s")
EOF
