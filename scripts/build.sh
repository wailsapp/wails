#!/usr/bin/env bash

echo "**** Checking if Wails passes unit tests ****"
if ! go test ./lib/... ./runtime/... ./cmd/...
then
    echo ""
    echo "ERROR: Unit tests failed!"
    exit 1;
fi

# Build runtime
echo "**** Building Runtime ****"
cd runtime/js
npm install
npm run build
cd ../..

cd cmd/wails
echo "**** Checking if Wails compiles ****"
if ! go build .
then
    echo ""
    echo "ERROR: Build failed!"
    exit 1;
fi

echo "**** Installing Wails locally ****"
if ! go install
then
    echo ""
    echo "ERROR: Install failed!"
    exit 1;
fi
cd ../..

echo "**** Tidying the mods! ****"
go mod tidy

echo "**** WE ARE DONE! ****"
