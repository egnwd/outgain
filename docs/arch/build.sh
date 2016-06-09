#!/usr/bin/env bash

set -eux

dot -Tpdf -o diagram.pdf diagram.dot
pdflatex architecture.tex

