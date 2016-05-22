#!/usr/bin/env bash
set -eux

(cd client; gulp)
(cd server; go build -v)
