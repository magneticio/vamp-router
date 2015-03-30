#!/bin/bash
set -e

# Set the version
VERSION=$1
if [ -z $VERSION ]; then
    echo "Enter a version"
    exit 1
fi

# clean up previous packages
rm -rf ./target/dist/
mkdir -p ./target/dist/



# build the app for linux/i386 and create a zip with necessary artifacts
GOOS=linux GOARCH=386 go build
cp vamp-router ./target/dist/
cp -r ./configuration ./target/dist/
cp -r ./examples ./target/dist/
cd ./target/dist/
zip -r vamp-router_${VERSION}_linux_386.zip *

 

