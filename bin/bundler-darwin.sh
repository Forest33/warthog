#!/bin/sh

version=$(cat version)

#go install github.com/asticode/go-astilectron-bundler/astilectron-bundler

mkdir -p distr

mkdir -p ../deploy/app/resources
cp -R ../resources/* ../deploy/app/resources
mv ../deploy/app/bind.go ../deploy/app/bind.go.tmp

cd ../deploy/app
astilectron-bundler -c ../../bin/bundler-darwin.json -ldflags X:main.UseBootstrap=true -ldflags X:main.AppVersion=${version} -ldflags "-s -w"

cp resources/icons/tray24.png ../../bin/distr/darwin-amd64/warthog.app/Contents/Resources/
rm -R resources
mv bind.go.tmp bind.go
rm bind_darwin_amd64.go

cd ../../bin/distr/darwin-amd64/warthog.app/Contents/MacOS/
upx -9 Warthog


