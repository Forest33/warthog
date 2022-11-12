#!/bin/sh

version=$(cat version)

mkdir -p distr

mkdir -p ../deploy/app/resources
cp -R ../resources/* ../deploy/app/resources
mv ../deploy/app/bind.go ../deploy/app/bind.go.tmp

cd ../deploy/app || exit
astilectron-bundler -c ../../bin/bundler-linux.json -ldflags X:main.UseBootstrap=true -ldflags X:main.AppVersion="${version}" -ldflags "-s -w"

rm -R resources
mv bind.go.tmp bind.go
rm bind_linux_amd64.go

cd ../../bin/distr/linux-amd64 || exit
upx -9 Warthog