#!/usr/bin/env bash

GITHUB_REF=${GITHUB_REF:-refs/heads/master}

git rev-parse $GITHUB_REF > /dev/null 2>&1
if [ $? != 0 ]
then
      echo "\$GITHUB_REF is invalid"
      exit 1
fi

if [[ $GITHUB_REF = refs/tags/* ]]
then
      # previous tag
      export PREVIOUS_REF=$(git tag --sort=refname --format='%(refname)' | grep $GITHUB_REF -B 1 | grep -m1 "")
else
      # latest tag
      export PREVIOUS_REF=$(git tag --sort=-refname --format='%(refname)' | grep -m1 "")
fi

cat <<EOF
## Changes
$(git log $GITHUB_REF...$PREVIOUS_REF --format="- [%h](../../commit/%h) %s")
EOF
