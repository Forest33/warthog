#!/bin/sh

version=$(cat version)

#sudo apt-get install gcc-multilib
#sudo apt-get install gcc-mingw-w64
#go install github.com/asticode/go-astilectron-bundler/astilectron-bundler

mkdir -p ../deploy/app/resources
cp -R ../resources/* ../deploy/app/resources
mv ../deploy/app/bind.go ../deploy/app/bind.go.tmp

cd ../deploy/app || exit
astilectron-bundler -c ../../bin/bundler-windows64.json -ldflags X:main.UseBootstrap=true -ldflags X:main.AppVersion="${version}" -ldflags "-s -w"

rm -R resources
mv bind.go.tmp bind.go
rm bind_windows_amd64.go
rm windows.syso

cd ../../bin/distr/windows-amd64 || exit
upx -9 Warthog.exe

