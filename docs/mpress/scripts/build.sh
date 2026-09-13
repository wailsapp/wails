#!/usr/bin/env bash
set -euo pipefail

# Cloudflare and CI build the same pinned release, without Go or Node.js.
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
build_tmp=$(mktemp -d)
trap 'rm -rf "$build_tmp"' EXIT
release_url="https://github.com/leaanthony/mpress/releases/download/v1.0.3"
archive="mpress-linux-amd64.tar.gz"
expected="97c3db3163789a8551e72b5929213d9fdb13ddf9a04f3cddbfced0484c01632a"
curl --fail --location --silent --show-error --retry 3 "$release_url/$archive" -o "$build_tmp/$archive"
echo "$expected  $build_tmp/$archive" | sha256sum --check --status
tar -xzf "$build_tmp/$archive" -C "$build_tmp"
test "$("$build_tmp/mpress" version)" = "mpress 1.0.3"
cd "$project_root"
# The imported home-page animation adds CSS classes at runtime.
"$build_tmp/mpress" build --strict --no-purge-css
"$build_tmp/mpress" check
python3 -m unittest discover -s docs/mpress/scripts -p 'test_*.py'
python3 docs/mpress/scripts/check_site.py docs/mpress/site
