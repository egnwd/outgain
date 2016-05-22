#!/usr/bin/env bash
set -eu

: "${HEROKU_OAUTH_TOKEN:=}"

if [[ $CIRCLE_PROJECT_USERNAME == "egnwd" &&
      $CIRCLE_BRANCH == "master" ]]; then
    APPNAME="outgain"
else
    APPNAME="outgain-$CIRCLE_PROJECT_USERNAME"
fi

if [[ -n $HEROKU_OAUTH_TOKEN ]]; then
    echo "Building slug ..."
    ./build_slug.sh app

    echo "Creating slug archive ..."
    tar czvf slug.tgz ./app

    echo "Deploying slug to $APPNAME ..."
    ./deploy.rb "$APPNAME" slug.tgz
else
    echo "WARNING: HEROKU_OAUTH_TOKEN missing, skipping deploy"
fi
