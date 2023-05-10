#!/usr/bin/bash
# See LICENSE.txt for copyright and licensing information about this file.

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";
PARENT="${SCRIPT_DIR%/*}"

cd "${PARENT}"
go build -cover -o ./bin/goglob ./cmd
export GOCOVERDIR=coverage/int
./bin/goglob "a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*" "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" > /dev/null
rm ./coverage/unit/*
go test -cover . -args -test.gocoverdir=./coverage/unit > /dev/null

echo Unit Test Coverage
go tool covdata textfmt -i=./coverage/unit  -o ./coverage/profile
go tool cover -func ./coverage/profile
