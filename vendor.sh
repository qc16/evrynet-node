#!/usr/bin/env bash
set -euo pipefail
# download vendor
go mod vendor
# clone project with .h files
# https://github.com/golang/go/issues/26366
pushd vendor/github.com/karalabe
rm -rf usb
git clone https://github.com/karalabe/usb.git
popd
chmod -R u+w vendor/github.com/karalabe/usb
