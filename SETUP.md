# Setup Guide for terraform-provider-coolify Fork

This guide explains how to set up and use your fork of terraform-provider-coolify.

## Quick Start: Using the Provider Locally

The **fastest way** to start using this provider is to build it locally and use Terraform's development overrides.

### 1. Build and Install

```bash
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify
make install
```

This builds the provider and installs it to `$(go env GOPATH)/bin` (typically `~/go/bin`).

### 2. Configure Terraform Development Override

Create or edit `~/.terraformrc` (macOS/Linux) or `%APPDATA%/terraform.rc` (Windows):

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/patrikwm/coolify" = "/Users/patrik/go/bin"
  }

  # For all other providers, use the standard registry
  direct {}
}
```

**Important:** Replace `/Users/patrik/go/bin` with your actual Go binary path:
```bash
echo $(go env GOPATH)/bin
```

### 3. Use in Your Terraform Configuration

Create a test directory:

```bash
mkdir ~/coolify-terraform-test
cd ~/coolify-terraform-test
```

Create `main.tf`:

```hcl
terraform {
  required_providers {
    coolify = {
      source = "patrikwm/coolify"
      # No version needed with dev_overrides
    }
  }
}

provider "coolify" {
  endpoint = "https://your-coolify.example.com/api/v1"
  # token = "..." # Or set COOLIFY_TOKEN env var
}

# Example: Create a new environment
resource "coolify_environment" "staging" {
  project_uuid = "your-project-uuid"
  name         = "staging"
}

# Example: Deploy a service with optional destination
resource "coolify_service" "app" {
  name             = "my-app"
  server_uuid      = "your-server-uuid"
  project_uuid     = "your-project-uuid"
  environment_name = coolify_environment.staging.name
  # destination_uuid is optional - auto-selects if server has only one destination
  instant_deploy   = true

  compose = <<EOF
services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
EOF
}

# Output the computed status
output "service_status" {
  value = coolify_service.app.status
}

output "server_status" {
  value = coolify_service.app.server_status
}
```

### 4. Run Terraform

```bash
export COOLIFY_TOKEN="your-api-token"
terraform init
terraform plan
terraform apply
```

**Note:** With dev_overrides, Terraform will warn that it's using an overridden provider. This is expected.

## Publishing Your Fork to GitHub

If you want to share your fork or use it from multiple machines:

### 1. Push to Your GitHub Repository

```bash
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify
git add .
git commit -m "feat: add environment resource and service enhancements"
git push origin main
```

### 2. Create a Release

#### Option A: Via GitHub UI

1. Go to https://github.com/patrikwm/terraform-provider-coolify
2. Click "Releases" → "Create a new release"
3. Create a new tag: `v0.1.0` (follow [Semantic Versioning](https://semver.org/))
4. Title: `v0.1.0 - Initial Fork Release`
5. Description: List your changes
6. Click "Publish release"

#### Option B: Via Command Line

```bash
git tag -a v0.1.0 -m "Initial fork release with environment resource"
git push origin v0.1.0
```

### 3. GitHub Actions Will Build Binaries

The `.github/workflows/release.yml` workflow will:
- Build binaries for all platforms (Linux, macOS, Windows, FreeBSD × amd64/arm64/386/arm)
- Create checksums
- **Optionally** sign with GPG (requires secrets - see below)
- Create a draft release

### 4. GPG Signing (Optional but Recommended for Terraform Registry)

If you plan to publish to the Terraform Registry, you **must** sign releases:

#### Generate GPG Key

```bash
# Generate key
gpg --full-generate-key
# Select: RSA and RSA, 4096 bits
# Use your GitHub email
# Set a strong passphrase

# Get fingerprint
gpg --list-secret-keys --keyid-format=long
# Look for the line like: sec   rsa4096/ABCD1234EFGH5678 2026-03-01

# Export public key (for Terraform Registry)
gpg --armor --export your-email@example.com > public.key

# Export private key (for GitHub Actions)
gpg --armor --export-secret-keys your-email@example.com > private.key
```

#### Add Secrets to GitHub

1. Go to https://github.com/patrikwm/terraform-provider-coolify/settings/secrets/actions
2. Add three secrets:
   - `GPG_PRIVATE_KEY`: Contents of `private.key` file
   - `PASSPHRASE`: Your GPG key passphrase
   - `GPG_FINGERPRINT`: Your key fingerprint (e.g., `ABCD1234EFGH5678`)

**Security:** Delete the `private.key` file after uploading to GitHub.

#### Using Without GPG

If you're not publishing to Terraform Registry, you can skip GPG signing. Comment out the signing section in `.goreleaser.yml`:

```yaml
# signs:
#   - artifacts: checksum
#     args: ...
```

### 5. Use the Published Release

Create `~/.terraformrc`:

```hcl
provider_installation {
  # Remove dev_overrides when using published releases

  filesystem_mirror {
    path    = "/tmp/terraform-plugins"
    include = ["patrikwm/coolify"]
  }

  direct {}
}
```

Download the release manually:

```bash
mkdir -p /tmp/terraform-plugins/registry.terraform.io/patrikwm/coolify/0.1.0/darwin_arm64
cd /tmp/terraform-plugins/registry.terraform.io/patrikwm/coolify/0.1.0/darwin_arm64
# Download from GitHub releases
# Adjust platform for your OS/arch
```

Or use the provider directly from GitHub (requires manual download):

```hcl
terraform {
  required_providers {
    coolify = {
      source  = "patrikwm/coolify"
      version = "~> 0.1.0"
    }
  }
}
```

## Publishing to Terraform Registry (Advanced)

To make your provider available via `registry.terraform.io`:

1. **Create a Terraform Registry Account**: https://registry.terraform.io/
2. **Verify GitHub OAuth**: Connect your GitHub account
3. **Add GPG Public Key**: Upload your `public.key`
4. **Publish Provider**: Follow https://developer.hashicorp.com/terraform/registry/providers/publishing
5. **Requirements**:
   - Repository must be named `terraform-provider-coolify`
   - Must have signed releases
   - Must follow Terraform provider conventions

After publishing, users can use:

```hcl
terraform {
  required_providers {
    coolify = {
      source  = "patrikwm/coolify"
      version = "~> 0.1.0"
    }
  }
}
```

And Terraform will automatically download from the registry.

## Maintenance

### Updating from Upstream

To merge updates from the original SierraJC repository:

```bash
# Add upstream remote (one-time)
git remote add upstream https://github.com/SierraJC/terraform-provider-coolify.git

# Fetch upstream changes
git fetch upstream

# Merge upstream main into your fork
git checkout main
git merge upstream/main

# Resolve conflicts if any
# Test and commit
git push origin main
```

### Updating OpenAPI Schema

```bash
# Fetch latest Coolify API schema
make fetch-schema

# Review changes
git diff tools/openapi.yml

# Regenerate code
make generate

# Test
make test

# Commit
git add .
git commit -m "chore: update to latest Coolify API schema"
```

## Environment Variables

For development and testing:

```bash
# Required for provider
export COOLIFY_ENDPOINT="https://your-coolify.example.com/api/v1"
export COOLIFY_TOKEN="your-api-token"

# Optional for development
export TF_LOG=DEBUG              # Verbose Terraform logging
export TF_ACC=1                  # Enable acceptance tests
export GOPATH="$HOME/go"         # Go workspace (usually default)
```

Create a `.env` file (ignored by git):

```bash
cp .env.example .env
# Edit .env with your values
```

Source it:
```bash
source .env  # or: set -a; source .env; set +a
```

## Troubleshooting

### "Provider not found" error

**Problem:** Terraform can't find the provider.

**Solution:** Check your `~/.terraformrc` file:
- Verify `dev_overrides` path matches `$(go env GOPATH)/bin`
- Ensure `make install` completed successfully
- Restart your terminal to reload rc files

### "Dev overrides in use" warning

**This is normal** when using local builds. Terraform warns you that it's using your local version instead of downloading from the registry.

### Provider builds but changes don't appear

**Solution:** Rebuild and reinstall:
```bash
make install
# In your Terraform project:
rm -rf .terraform .terraform.lock.hcl
terraform init
```

### Acceptance tests fail

**Common causes:**
1. `.env` file not configured
2. Coolify instance not reachable
3. API token invalid or expired
4. Required resources don't exist (project, server, etc.)

**Solution:**
```bash
# Verify connection
curl -H "Authorization: Bearer $COOLIFY_TOKEN" "$COOLIFY_ENDPOINT/../health"

# Should return: OK
```

## Next Steps

1. **Read the docs**: Check `docs/` for resource documentation
2. **Review examples**: See `examples/` for usage patterns
3. **Contribute**: See [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow
4. **Report issues**: https://github.com/patrikwm/terraform-provider-coolify/issues

Happy Terraforming! 🚀
