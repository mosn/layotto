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
go test -count=1 -failfast -timeout 120s ./... -coverprofile cover.out -covermode=atomic
cat cover.out >> ../../cover.out
cd ../..

#go tool cover -html=cover.out -o coverage.html
#go tool cover -func=cover.out -o func.out
#echo test func-coverage $(tail -1 func.out | awk '{print $3}')

#for d in $(go list ./pkg/...); do
#    echo "--------Run test package: $d"
#    GO111MODULE=off go test -gcflags="all=-N -l" -v -coverprofile=profile.out -covermode=atomic $d
#    echo "--------Finish test package: $d"
#    if [ -f profile.out ]; then
#        cat profile.out >> coverage.txt
#        rm profile.out
#    fi
#done
#
#cd components
#for d in $(go list ./...); do
#    echo "--------Run test package: $d"
#    GO111MODULE=off go test -gcflags=-l -v -coverprofile=profile.out -covermode=atomic $d
#    echo "--------Finish test package: $d"
#    if [ -f profile.out ]; then
#	cat profile.out >> coverage.txt
#	rm profile.out
#    fi
#done
