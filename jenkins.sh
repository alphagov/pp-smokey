#!/bin/sh

set -x

# This removes rbenv shims from the PATH where there is no
# .ruby-version file. This is because certain gems call their
# respective tasks with ruby -S which causes the following error to
# appear: ruby: no Ruby script found in input (LoadError).
if [ ! -f .ruby-version ]; then
  export PATH=$(printf $PATH | awk 'BEGIN { RS=":"; ORS=":" } !/rbenv/' | sed 's/:$//')
fi

bundle install --path "${HOME}/bundles/${JOB_NAME}" --deployment --quiet

if [ "$PP_EXCLUDE_PATTERN" != "" ]; then
  EXCLUDE_FEATURE="-e '${PP_EXCLUDE_PATTERN}'"
fi

bundle exec cucumber --format json --out ${WORKSPACE}/results.json $EXCLUDE_FEATURE
