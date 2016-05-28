#!/usr/bin/env bash
set -eux

./build_mruby.sh
(cd client; npm install && gulp)
(cd server; go get -t -d -v ./... && go build -v)
