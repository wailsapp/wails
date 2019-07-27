#!/usr/bin/env bash

# Build runtime
echo "**** Building Runtime ****"
cd runtime/js
npm run build
cd ../..

echo "**** Packing Assets ****"
mewn

echo "**** Installing Wails locally ****"
cd cmd/wails
go install
cd ../..

echo "**** Tidying the mods! ****"
go mod tidy

echo "**** WE ARE DONE! ****"
