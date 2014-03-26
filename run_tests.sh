#!/bin/bash -e

echo
echo 'Note: you need rbenv and you may need to run:'
echo '> rbenv install $(cat .ruby-version)'
echo

bundle install
bundle exec rake
