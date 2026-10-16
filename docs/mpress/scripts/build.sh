#!/usr/bin/env bash
set -euo pipefail

# Cloudflare and CI download the latest release, without Go or Node.js.
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
build_tmp=$(mktemp -d)
trap 'rm -rf "$build_tmp"' EXIT
# Resolve latest once so the archive and checksums belong to the same release.
release_page=$(curl --fail --location --silent --show-error --retry 3 --head \
  --output /dev/null --write-out '%{url_effective}' \
  https://github.com/leaanthony/mpress/releases/latest)
release_tag="${release_page##*/}"
release_url="https://github.com/leaanthony/mpress/releases/download/$release_tag"
archive="mpress-linux-amd64.tar.gz"
echo "Downloading M-Press $release_tag"
curl --fail --location --silent --show-error --retry 3 "$release_url/$archive" -o "$build_tmp/$archive"
curl --fail --location --silent --show-error --retry 3 "$release_url/checksums.txt" -o "$build_tmp/checksums.txt"
awk -v archive="$archive" '$2 == archive {print}' "$build_tmp/checksums.txt" > "$build_tmp/$archive.sha256"
(cd "$build_tmp" && sha256sum --check --strict "$archive.sha256")
tar -xzf "$build_tmp/$archive" -C "$build_tmp"
test "$("$build_tmp/mpress" version)" = "mpress ${release_tag#v}"
cd "$project_root"
python3 docs/mpress/scripts/check_translations.py --mpress "$build_tmp/mpress"
# The imported home-page animation adds CSS classes at runtime.
"$build_tmp/mpress" build --strict --no-purge-css
"$build_tmp/mpress" check
python3 -m unittest discover -s docs/mpress/scripts -p 'test_*.py'
python3 docs/mpress/scripts/check_site.py docs/mpress/site
