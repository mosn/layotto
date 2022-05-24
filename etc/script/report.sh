#!/usr/bin/env bash

set -e
echo "" > cover.out
echo "test pkg"
go test -count=1 -failfast -timeout 120s ./pkg/... -coverprofile cover.out -covermode=atomic


cd components
echo "" > cover.out
echo "test components"
go test -count=1 -failfast -timeout 120s ./... -coverprofile cover.out -covermode=atomic
cat cover.out >> ../cover.out
cd ..


cd sdk/go-sdk
echo "" > cover.out
echo "test go-sdk"
go test -count=1 -failfast -timeout 120s $(go list ./... | grep -v runtime) -coverprofile cover.out -covermode=atomic
cat cover.out >> ../../cover.out
cd ../..
