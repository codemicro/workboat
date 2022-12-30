#!/usr/bin/env bash

set -e

REPO_URL=$1
AUTH_STR=$2
COMMAND=$3

GIT_TERMINAL_PROMPTS=0

if [[ "$REPO_URL" == "" || "$COMMAND" == "" ]]; then
  exit 3
fi

if [[ $AUTH_STR != "" ]]; then
  echo "$AUTH_STR" > gitcred.txt
  git config --global credential.helper "store --file $(pwd)/gitcred.txt"
fi

git clone $REPO_URL repository
cd repository

$COMMAND