#!/bin/sh

version=$(cat ../../version)

cd ../../distr/linux-amd64
tar czvf ../../installer/linux/Warthog-${version}-linux-x86-64.tar.gz *

