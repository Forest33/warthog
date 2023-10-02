#!/usr/bin/env sh

version=$(cat ../../version)

for a in "x86-64,amd64,Intel" "arm64,arm64,M1"; do
  IFS=","
  set -- $a
  arch=$1; shift
  path=$1; shift
  name=$1; shift

  create-dmg \
    --volname "Warthog ${name} Installer" \
    --volicon "app.icns" \
    --window-pos 200 120 \
    --window-size 800 400 \
    --icon-size 100 \
    --icon "Warthog.app" 200 190 \
    --hide-extension "Warthog.app" \
    --app-drop-link 600 185 \
    "Warthog-${version}-darwin-${arch}.dmg" \
    "../../distr/darwin-${path}/"
done

