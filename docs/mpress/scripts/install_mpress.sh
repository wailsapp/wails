#!/usr/bin/env bash
set -euo pipefail

# Pin release builds and retrospective snapshots to the same renderer.
version=${MPRESS_VERSION:-v1.0.25}
[[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo 'Invalid M-Press version' >&2; exit 1; }
destination=${1:?usage: install_mpress.sh DIRECTORY}
mkdir -p "$destination"
archive=mpress-linux-amd64.tar.gz
release_url="https://github.com/leaanthony/mpress/releases/download/$version"
curl --fail --location --silent --show-error --retry 3 "$release_url/$archive" -o "$destination/$archive"
curl --fail --location --silent --show-error --retry 3 "$release_url/checksums.txt" -o "$destination/checksums.txt"
awk -v archive="$archive" '$2 == archive {print}' "$destination/checksums.txt" > "$destination/$archive.sha256"
(cd "$destination" && sha256sum --check --strict "$archive.sha256")
tar -xzf "$destination/$archive" -C "$destination"
test "$("$destination/mpress" version)" = "mpress ${version#v}"
