# 🎯 Fork Setup Complete - Next Steps

Your terraform-provider-coolify fork is now production-ready! Here's what was updated and what you should do next.

## ✅ What Was Updated

### DocumentationFiles
1. **[README.md](README.md)**
   - Updated badges and links to point to `patrikwm/terraform-provider-coolify`
   - Added "Enhanced Features" section highlighting new capabilities
   - Updated resource support table with fork-specific enhancements
   - Added comprehensive usage instructions (local dev, GitHub releases, registry)
   - Removed codecov badge (you can add back if you set up codecov)

2. **[CONTRIBUTING.md](CONTRIBUTING.md)**
   - **NEW**: Complete contribution guide (was empty before)
   - Development setup instructions
   - Code generation workflow
   - Testing guidelines
   - How to add new resources
   - Commit conventions (Conventional Commits)
   - Release process with GPG signing instructions
   - Architecture overview

3. **[SETUP.md](SETUP.md)**
   - **NEW**: Comprehensive setup guide
   - Quick start for local development
   - Publishing releases to GitHub
   - GPG signing setup (optional)
   - Publishing to Terraform Registry (advanced)
   - Troubleshooting common issues
   - Environment variable reference

### CI/CD Workflows
4. **[.github/workflows/test.yml](.github/workflows/test.yml)**
   - Made codecov optional (won't fail if secrets not set)
   - Updated for patrikwm repository context

5. **[.github/workflows/release.yml](.github/workflows/release.yml)**
   - Already configured correctly
   - Requires GPG secrets for signing (optional)

6. **[.goreleaser.yml](.goreleaser.yml)**
   - Updated changelog footer to point to patrikwm repository

## 🚀 Recommended Next Steps

### 1. Start Using It Locally (Easiest)

```bash
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify

# Build and install
make install

# Configure Terraform to use local build
# Create ~/.terraformrc with:
provider_installation {
  dev_overrides {
    "registry.terraform.io/patrikwm/coolify" = "/Users/patrik/go/bin"
  }
  direct {}
}

# Start using in your Terraform configs!
```

See [SETUP.md](SETUP.md#quick-start-using-the-provider-locally) for detailed instructions.

### 2. Push to Your GitHub Repository

```bash
# Review changes
git status
git diff

# Commit all updates
git add .
git commit -m "docs: update for patrikwm fork with enhanced features

- Add coolify_environment resource
- Make destination_uuid optional on services
- Add computed status attributes to services
- Update documentation and CI/CD for fork
- Add comprehensive CONTRIBUTING and SETUP guides"

# Push to your repository
git push origin main
```

### 3. Create Your First Release (Optional)

```bash
# Tag the release
git tag -a v0.1.0 -m "Initial fork release with environment resource and service enhancements"
git push origin v0.1.0
```

The GitHub Action will:
- Build binaries for all platforms
- Create checksums
- Create a draft release (you review and publish)

**Note:** Without GPG signing, binaries will be published but won't be signed. This is fine for personal use but required for Terraform Registry.

### 4. Set Up GPG Signing (Recommended if Publishing)

Only needed if you want to:
- Publish to Terraform Registry
- Allow others to verify your releases

See [SETUP.md](SETUP.md#4-gpg-signing-optional-but-recommended-for-terraform-registry) for instructions.

## 📋 Repository Checklist

- ✅ Code is up-to-date with enhancements
- ✅ Documentation updated
- ✅ CI/CD workflows configured
- ⏳ Push to GitHub (you do this)
- ⏳ Create first release (optional)
- ⏳ Set up GPG signing (optional)
- ⏳ Publish to Terraform Registry (optional)

## 🛠️ Using the Provider

### Quick Example

```hcl
terraform {
  required_providers {
    coolify = {
      source = "patrikwm/coolify"
      # version not needed with dev_overrides
    }
  }
}

provider "coolify" {
  endpoint = "https://your-coolify.example.com/api/v1"
  # token via COOLIFY_TOKEN env var
}

# NEW: Environment resource (not in upstream!)
resource "coolify_environment" "staging" {
  project_uuid = "your-project-uuid"
  name         = "staging"
}

# ENHANCED: Optional destination_uuid
resource "coolify_service" "app" {
  name             = "my-app"
  server_uuid      = "your-server-uuid"
  project_uuid     = "your-project-uuid"
  environment_name = coolify_environment.staging.name
  # destination_uuid omitted - auto-selects!
  instant_deploy   = true

  compose = <<EOF
services:
  web:
    image: nginx:latest
EOF
}

# NEW: Computed status attributes
output "service_status" {
  value = coolify_service.app.status  # e.g., "running:healthy"
}
```

## 📚 Documentation Structure

- **[README.md](README.md)** - Overview, features, quick links
- **[SETUP.md](SETUP.md)** - Complete setup guide (local dev, releases, registry)
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Developer guide
- **[docs/](docs/)** - Auto-generated Terraform documentation
- **[examples/](examples/)** - Usage examples

## 🐛 Known Issues / Future Work

1. **Acceptance tests disabled** - CI comment says test server in maintenance
2. **Codecov not configured** - Removed from badges, but workflows still reference it (won't fail)
3. **Destinations data source** - Mentioned in plan but deferred (no dedicated API endpoint)

## 🔗 Useful Links

- **This Fork**: https://github.com/patrikwm/terraform-provider-coolify
- **Upstream**: https://github.com/SierraJC/terraform-provider-coolify
- **Coolify**: https://github.com/coollabsio/coolify
- **Coolify API Docs**: https://coolify.io/docs/api-reference/introduction
- **Terraform Registry**: https://registry.terraform.io/
- **Terraform Provider Framework**: https://developer.hashicorp.com/terraform/plugin/framework

## 🤝 Contribution Workflow

1. Fork is ready for contributions
2. Follow [CONTRIBUTING.md](CONTRIBUTING.md) guidelines
3. Use Conventional Commits
4. Run `make test` before committing
5. Create PR with clear description

## 💡 Tips

- **Local development**: Use dev_overrides in ~/.terraformrc
- **Quick rebuild**: `make install` after code changes
- **Update schema**: `make fetch-schema && make generate`
- **Test without Terraform**: `go test -v ./internal/service/...`
- **Debug**: Set `TF_LOG=DEBUG` for verbose Terraform output

---

**You're all set!** 🎉

Next: Run `make install` and configure your ~/.terraformrc to start using the provider locally.

For questions or issues, check [SETUP.md](SETUP.md) or [CONTRIBUTING.md](CONTRIBUTING.md).
