#!/usr/bin/env bash

set -eux

go get -d github.com/mitchellh/go-mruby
(cd $GOPATH/src/github.com/mitchellh/go-mruby; make)
cp $GOPATH/src/github.com/mitchellh/go-mruby/libmruby.a server
