#!/bin/sh
set -e

REPO="Dat-one-dev/Aksara"
BIN="aksara"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

os="$(uname -s)"
arch="$(uname -m)"

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  aarch64|arm64) arch="arm64" ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

case "$os" in
  Linux) os="Linux" ;;
  Darwin) os="Darwin" ;;
  *) echo "unsupported os: $os (windows users: grab Aksara_Windows_amd64.tar.gz from the latest release)" >&2; exit 1 ;;
esac

url="https://github.com/$REPO/releases/latest/download/Aksara_${os}_${arch}.tar.gz"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "downloading $url"
curl -fsSL "$url" -o "$tmp/aksara.tar.gz"
tar -xzf "$tmp/aksara.tar.gz" -C "$tmp"

mkdir -p "$INSTALL_DIR"
mv "$tmp/$BIN" "$INSTALL_DIR/$BIN"
chmod +x "$INSTALL_DIR/$BIN"

echo "installed to $INSTALL_DIR/$BIN"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "note: $INSTALL_DIR is not on your PATH, add it or move the binary somewhere it is" ;;
esac
