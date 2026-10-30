#!/usr/bin/env bash
set -euo pipefail

project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)
build_tmp=$(mktemp -d)
trap 'rm -rf "$build_tmp"' EXIT
if [[ -z "${MPRESS_BIN:-}" ]]; then
  bash "$project_root/docs/mpress/scripts/install_mpress.sh" "$build_tmp"
  export MPRESS_BIN="$build_tmp/mpress"
fi
cd "$project_root"
# Clean/shallow Cloudflare checkouts need release tags as well as HEAD.
git fetch --tags --no-recurse-submodules origin
python3 docs/mpress/scripts/check_translations.py --mpress "$MPRESS_BIN"
python3 docs/mpress/scripts/release_docs.py build --mpress "$MPRESS_BIN"
python3 -m unittest discover -s docs/mpress/scripts -p 'test_*.py'
