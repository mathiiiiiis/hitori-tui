#!/usr/bin/env bash
set -euo pipefail

repo="mathiiiiiis/hitori-tui"
install_dir="${HITORI_INSTALL_DIR:-$HOME/.local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

case "$os" in
  linux|darwin) ;;
  *) echo "unsupported OS: $os (use the .zip release on Windows)" >&2; exit 1 ;;
esac

version=$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" | grep -m1 '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
if [ -z "$version" ]; then
  echo "could not determine latest release" >&2
  exit 1
fi

archive="hitori-${version}-${os}-${arch}.tar.gz"
url="https://github.com/$repo/releases/download/v${version}/${archive}"

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

echo "downloading hitori v${version} for ${os}/${arch}..."
curl -fsSL "$url" -o "$tmpdir/$archive"
tar -xzf "$tmpdir/$archive" -C "$tmpdir"

mkdir -p "$install_dir"
mv "$tmpdir/hitori" "$install_dir/hitori"
chmod +x "$install_dir/hitori"

echo "installed hitori v${version} to $install_dir/hitori"
case ":$PATH:" in
  *":$install_dir:"*) ;;
  *) echo "note: $install_dir is not on your PATH" ;;
esac
