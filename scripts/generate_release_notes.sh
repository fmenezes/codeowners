#!/usr/bin/env bash

if [ -z "$GITHUB_REF" ]
then
      echo "\$GITHUB_REF is empty"
      exit 1
fi

if [[ ! $GITHUB_REF = refs/tags/* ]]
then
      echo "\$GITHUB_REF does not start with 'refs/tags/'"
      exit 1
fi

export PREVIOUS_TAG_REF=$(git tag --sort=refname --format='%(refname)' | grep $GITHUB_REF -B 1 | grep -m1 "")

cat <<EOF
## Changes
$(git log $GITHUB_REF...$PREVIOUS_TAG_REF --format="- [%h](../../commit/%h) %s")
EOF
