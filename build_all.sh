#!/usr/bin/env bash
set -eux

(cd client; npm install && gulp)
(cd server; go get -t -d -v ./... && go build -v)
(cd runner; cargo build --release ${CARGO_FLAGS:-})
