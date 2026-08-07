#!/bin/sh
set -eu

target=${1-}
draft_file=${2-}
if [ -z "$target" ]; then
  target='https://github.com/wailsapp/wails.git'
fi

run_contribution() {
  mpress_binary=$1
  if [ -n "$draft_file" ]; then
    exec "$mpress_binary" contribute --branch 'master' "$target" --draft-file "$draft_file"
  fi
  exec "$mpress_binary" contribute --branch 'master' "$target"
}

if command -v mpress >/dev/null 2>&1; then
  run_contribution mpress
fi

release_base=${MPRESS_RELEASE_BASE_URL:-https://github.com/leaanthony/mpress/releases/latest/download}
case "$release_base" in
  https://*) ;;
  *) echo "MPRESS_RELEASE_BASE_URL must use HTTPS." >&2; exit 1 ;;
esac

download() {
  source_url=$1
  destination=$2
  if command -v curl >/dev/null 2>&1; then
    if curl --proto '=https' --tlsv1.2 --fail --location --silent --show-error --retry 3 --retry-delay 1 --output "$destination" "$source_url"; then
      return 0
    fi
  fi
  if command -v wget >/dev/null 2>&1; then
    if wget --https-only --tries=3 --quiet --output-document="$destination" "$source_url"; then
      return 0
    fi
  fi
  echo "Could not download $source_url. Install curl or wget and try again." >&2
  return 1
}

system=$(uname -s | tr '[:upper:]' '[:lower:]')
machine=$(uname -m)
case "$system" in
  darwin) system=darwin ;;
  linux) system=linux ;;
  *) echo "M-Press does not provide a release for $(uname -s)." >&2; exit 1 ;;
esac
case "$machine" in
  x86_64|amd64) machine=amd64 ;;
  arm64|aarch64) machine=arm64 ;;
  *) echo "M-Press does not provide a release for architecture $machine." >&2; exit 1 ;;
esac

if ! command -v tar >/dev/null 2>&1; then
  echo "The tar command is required to unpack M-Press." >&2
  exit 1
fi

temporary_directory=$(mktemp -d 2>/dev/null || mktemp -d -t mpress)
cleanup() { rm -rf "$temporary_directory"; }
trap cleanup EXIT HUP INT TERM

asset="mpress-${system}-${machine}.tar.gz"
archive="$temporary_directory/$asset"
checksums="$temporary_directory/checksums.txt"
download "$release_base/$asset" "$archive"
download "$release_base/checksums.txt" "$checksums"

expected=$(awk -v asset="$asset" '$2 == asset || $2 == "*" asset { print $1; exit }' "$checksums")
if [ -z "$expected" ]; then
  echo "The M-Press release does not contain a checksum for $asset." >&2
  exit 1
fi
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$archive" | awk '{ print $1 }')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$archive" | awk '{ print $1 }')
else
  echo "A SHA-256 tool is required. Install sha256sum or shasum and try again." >&2
  exit 1
fi
if [ "$actual" != "$expected" ]; then
  echo "Checksum verification failed for $asset." >&2
  exit 1
fi

tar -xzf "$archive" -C "$temporary_directory"
binary="$temporary_directory/mpress"
if [ ! -f "$binary" ]; then
  echo "The M-Press release archive does not contain the mpress executable." >&2
  exit 1
fi
chmod +x "$binary"
run_contribution "$binary"
