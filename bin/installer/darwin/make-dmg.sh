#!/usr/bin/env sh

#genisoimage -V Warthog -D -R -apple -no-pad -o Warthog.dmg ../../distr/darwin-amd64/
create-dmg \
  --volname "Warthog Installer" \
  --volicon "app.icns" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "Warthog.app" 200 190 \
  --hide-extension "Warthog.app" \
  --app-drop-link 600 185 \
  "Warthog.dmg" \
  "../../distr/darwin-amd64/"
