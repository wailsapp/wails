#!/usr/bin/env bash
if [ "$#" != "1" ]; then 
  echo "Tag required"
  exit 1
fi
TAG=${1}
cat << EOF > cmd/version.go
package cmd

// Version - Wails version
const Version = "${TAG}"
EOF

# Build runtime
cd runtime/js
npm run build

cd ../..

git add cmd/version.go
git commit cmd/version.go -m "Bump to ${TAG}" 
git tag ${TAG}
