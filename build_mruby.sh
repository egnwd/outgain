#!/usr/bin/env bash

set -eux

go get -d github.com/mitchellh/go-mruby
GO_MRUBY=$(go list -f '{{.Dir}}' github.com/mitchellh/go-mruby)
(cd $GO_MRUBY; make)
cp $GO_MRUBY/libmruby.a server
