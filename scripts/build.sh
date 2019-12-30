#!/usr/bin/env bash

# Build runtime
echo "**** Building Runtime ****"
cd runtime/js
npm install
npm run build
cd ../..

echo "**** Packing Assets ****"
cd cmd
mewn
cd ..
cd lib/renderer
mewn
cd ../..

echo "**** Installing Wails locally ****"
cd cmd/wails
go install
cd ../..

echo "**** Tidying the mods! ****"
go mod tidy

echo "**** WE ARE DONE! ****"
