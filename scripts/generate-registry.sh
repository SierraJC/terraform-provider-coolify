#!/usr/bin/env bash
# Generate Tofu/Terraform HTTP Registry Protocol JSON for all GH Releases.
#
# Reads release assets from GitHub, generates /docs/v1/providers/wappsdev/coolify/...
# manifests. Intended to be idempotent — running twice produces same output.
#
# Requires: gh CLI, jq, curl, gpg. Auth via gh auth login (or GH_TOKEN env).
#
# Usage:
#   ./scripts/generate-registry.sh                      # all releases
#   ./scripts/generate-registry.sh v1.0.0               # one release

set -euo pipefail

REPO="wappsdev/terraform-provider-coolify"
NAMESPACE="wappsdev"
TYPE="coolify"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DOCS_ROOT="$REPO_ROOT/docs"
REGISTRY_ROOT="$DOCS_ROOT/v1/providers/$NAMESPACE/$TYPE"

GPG_KEY_FILE="$DOCS_ROOT/wappsdev-sigs.gpg"
if [ ! -f "$GPG_KEY_FILE" ]; then
  echo "ERROR: $GPG_KEY_FILE missing — run Task 3 first" >&2
  exit 1
fi
GPG_ASCII=$(cat "$GPG_KEY_FILE")
GPG_FPR=$(gpg --show-keys --with-colons "$GPG_KEY_FILE" | awk -F: '/^fpr:/ {print $10; exit}')
GPG_KEY_ID="${GPG_FPR: -16}"

if [ -z "$GPG_KEY_ID" ]; then
  echo "ERROR: could not extract GPG key ID from $GPG_KEY_FILE" >&2
  exit 1
fi

if [ "${1:-}" = "" ]; then
  RELEASES=$(gh release list --repo "$REPO" --limit 100 --json tagName --jq '.[].tagName')
else
  RELEASES="$1"
fi

for RELEASE_TAG in $RELEASES; do
  VERSION="${RELEASE_TAG#v}"
  echo ">>> Processing $RELEASE_TAG (version $VERSION)"

  SHASUMS_TMP=$(mktemp)
  gh release download "$RELEASE_TAG" --repo "$REPO" \
    -p "*SHA256SUMS" -O "$SHASUMS_TMP" --clobber

  while IFS= read -r line; do
    SHASUM=$(echo "$line" | awk '{print $1}')
    ASSET=$(echo "$line" | awk '{print $2}' | sed 's|^\*||')

    case "$ASSET" in
      *_SHA256SUMS|*_manifest.json) continue ;;
    esac

    # Extract os + arch from name: terraform-provider-coolify_1.0.0_darwin_arm64.zip
    PLATFORM_PART=$(echo "$ASSET" | sed -E "s|^terraform-provider-coolify_${VERSION}_||; s|\.zip$||")
    OS="${PLATFORM_PART%_*}"
    ARCH="${PLATFORM_PART##*_}"

    OUTPUT_DIR="$REGISTRY_ROOT/$VERSION/download/$OS"
    mkdir -p "$OUTPUT_DIR"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$RELEASE_TAG/$ASSET"
    SHASUMS_URL="https://github.com/$REPO/releases/download/$RELEASE_TAG/terraform-provider-coolify_${VERSION}_SHA256SUMS"
    SHASUMS_SIG_URL="https://github.com/$REPO/releases/download/$RELEASE_TAG/terraform-provider-coolify_${VERSION}_SHA256SUMS.sig"

    GPG_ASCII_ESCAPED=$(echo "$GPG_ASCII" | jq -Rsr @json)

    cat > "$OUTPUT_DIR/$ARCH" <<JSON_EOF
{
  "protocols": ["6.0"],
  "os": "$OS",
  "arch": "$ARCH",
  "filename": "$ASSET",
  "download_url": "$DOWNLOAD_URL",
  "shasums_url": "$SHASUMS_URL",
  "shasums_signature_url": "$SHASUMS_SIG_URL",
  "shasum": "$SHASUM",
  "signing_keys": {
    "gpg_public_keys": [
      {
        "key_id": "$GPG_KEY_ID",
        "ascii_armor": $GPG_ASCII_ESCAPED,
        "trust_signature": "",
        "source": "wappsdev",
        "source_url": "https://registry.wapps.co/wappsdev-sigs.gpg"
      }
    ]
  }
}
JSON_EOF
    echo "  wrote $OUTPUT_DIR/$ARCH"
  done < "$SHASUMS_TMP"
  rm "$SHASUMS_TMP"
done

echo ">>> Generating versions endpoint"
VERSIONS_JSON=$(jq -n '[]')

for VERSION_DIR in "$REGISTRY_ROOT"/*/; do
  [ -d "$VERSION_DIR" ] || continue
  VERSION=$(basename "$VERSION_DIR")
  PLATFORMS=$(jq -n '[]')
  for OS_DIR in "$VERSION_DIR/download"/*/; do
    [ -d "$OS_DIR" ] || continue
    OS=$(basename "$OS_DIR")
    for ARCH_FILE in "$OS_DIR"*; do
      [ -f "$ARCH_FILE" ] || continue
      ARCH=$(basename "$ARCH_FILE")
      PLATFORMS=$(echo "$PLATFORMS" | jq --arg os "$OS" --arg arch "$ARCH" '. + [{"os":$os,"arch":$arch}]')
    done
  done
  VERSIONS_JSON=$(echo "$VERSIONS_JSON" | jq --arg v "$VERSION" --argjson p "$PLATFORMS" '. + [{"version":$v,"protocols":["6.0"],"platforms":$p}]')
done

echo "$VERSIONS_JSON" | jq '{versions: .}' > "$REGISTRY_ROOT/versions"
echo "  wrote $REGISTRY_ROOT/versions"
echo ">>> Done."
