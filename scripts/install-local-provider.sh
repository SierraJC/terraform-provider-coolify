#!/bin/bash

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROVIDER_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
PLUGIN_DIR="$HOME/.terraform.d/plugins/registry.terraform.io/sierrajc/coolify/dev/${OS}_${ARCH}"

echo -e "${BLUE}🔨 Building Terraform Coolify provider...${NC}"
echo "   Directory: $PROVIDER_DIR"

if [ -f "$PROVIDER_DIR/tools/tfplugingen-openapi.yml" ]; then
    echo "📝 Generating code from OpenAPI..."
    (cd "$PROVIDER_DIR" && make generate) || echo "⚠️  Generation failed, continuing..."
fi

echo "📦 Installing provider to $PLUGIN_DIR..."
mkdir -p "$PLUGIN_DIR"

BINARY_NAME="terraform-provider-coolify_dev_${OS}_${ARCH}"
BINARY_PATH="$PLUGIN_DIR/$BINARY_NAME"

echo "🔨 Building provider..."
(cd "$PROVIDER_DIR" && go build -o "$BINARY_PATH" .)

if [ ! -f "$BINARY_PATH" ]; then
    echo "❌ Build failed"
    exit 1
fi

chmod +x "$PLUGIN_DIR/$BINARY_NAME"

echo ""
echo -e "${GREEN}✅ Provider installed successfully!${NC}"
echo ""
echo -e "${BLUE}Location:${NC} $PLUGIN_DIR/$BINARY_NAME"
echo ""
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. In your Terraform project, create a file with:"
echo ""
echo "   terraform {"
echo "     required_providers {"
echo "       coolify = {"
echo "         source  = \"registry.terraform.io/sierrajc/coolify\""
echo "         version = \"dev\""
echo "       }"
echo "     }"
echo "   }"
echo ""
echo "2. Configure the provider:"
echo ""
echo "   provider \"coolify\" {"
echo "     # token will be read from COOLIFY_TOKEN"
echo "   }"
echo ""
echo "3. Initialize Terraform:"
echo ""
echo "   terraform init"
echo ""
echo -e "${YELLOW}💡 Tip:${NC} Export your API token:"
echo "   export COOLIFY_TOKEN=\"your-api-token\""
echo "   export COOLIFY_ENDPOINT=\"https://app.coolify.io/api/v1\"  # Optional"

