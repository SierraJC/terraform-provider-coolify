# 🚀 Quick Start Checklist

Copy this file: `cp QUICKSTART.md ~/Desktop/` for easy reference!

## ✅ Immediate Next Steps (Choose One Path)

### Path A: Use Locally (Fastest - 5 minutes)

```bash
# 1. Build and install
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify
make install

# 2. Set up Terraform override
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "registry.terraform.io/patrikwm/coolify" = "/Users/patrik/go/bin"
  }
  direct {}
}
EOF

# 3. Test it works
mkdir -p ~/coolify-test && cd ~/coolify-test
cat > test.tf <<'EOF'
terraform {
  required_providers {
    coolify = {
      source = "patrikwm/coolify"
    }
  }
}

provider "coolify" {
  endpoint = "YOUR_COOLIFY_URL/api/v1"
  token    = "YOUR_API_TOKEN"
}

# Example: List all servers
data "coolify_servers" "all" {}

output "servers" {
  value = data.coolify_servers.all
}
EOF

export COOLIFY_TOKEN="your-token-here"
terraform init
terraform plan
```

**Done!** You're using the enhanced provider with:
- ✨ `coolify_environment` resource
- ✨ Optional `destination_uuid` on services
- ✨ Computed `status` and `server_status` attributes

---

### Path B: Publish to GitHub (For Sharing - 10 minutes)

```bash
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify

# 1. Review what changed
git status
git diff README.md

# 2. Commit everything
git add .
git commit -m "docs: setup fork with enhancements

- Add environment resource
- Make destination_uuid optional
- Add status attributes to services
- Complete documentation overhaul"

# 3. Push to your repo
git push origin main

# 4. Create first release
git tag -a v0.1.0 -m "Initial fork release"
git push origin v0.1.0

# GitHub Actions will build binaries automatically!
# Check: https://github.com/patrikwm/terraform-provider-coolify/actions
```

Then other people (or you from other machines) can download from GitHub releases.

---

## 📋 Optional Enhancements

### Set Up GPG Signing (For Terraform Registry)

```bash
# Generate key
gpg --full-generate-key
# Choose: RSA, 4096 bits, use your GitHub email

# Get fingerprint
gpg --list-secret-keys --keyid-format=long
# Note the fingerprint (e.g., ABCD1234EFGH5678)

# Export for GitHub Actions
gpg --armor --export YOUR_EMAIL > public.key
gpg --armor --export-secret-keys YOUR_EMAIL > private.key

# Add to GitHub Secrets:
# https://github.com/patrikwm/terraform-provider-coolify/settings/secrets/actions
# - GPG_PRIVATE_KEY: paste contents of private.key
# - PASSPHRASE: your GPG passphrase
# - GPG_FINGERPRINT: your fingerprint

# Delete private key file
rm private.key
```

### Set Up Environment File

```bash
cd /Users/patrik/src/external/github-forks/coolify/open_tofu/terraform-provider-coolify
cp .env.example .env
# Edit .env with your Coolify details
nano .env
```

---

## 📚 Documentation Reference

| File | Purpose |
|------|---------|
| [FORK_SETUP.md](FORK_SETUP.md) | ⭐ **START HERE** - Complete overview of changes |
| [SETUP.md](SETUP.md) | Detailed setup instructions for all use cases |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Developer guide for contributing |
| [README.md](README.md) | Project overview and quick links |
| [docs/](docs/) | Auto-generated Terraform docs |
| [examples/](examples/) | Usage examples |

---

## 🎯 What's Different in This Fork?

Compared to upstream SierraJC/terraform-provider-coolify:

| Feature | Upstream | This Fork |
|---------|----------|-----------|
| `coolify_environment` resource | ❌ Blocked | ✅ Full CRUD support |
| Service `destination_uuid` | Required (empty string default) | ✅ Optional (auto-selects) |
| Service status attributes | ❌ Discarded | ✅ `status`, `server_status` |
| OpenAPI schemas | Incomplete Environment | ✅ Fixed with UUID |
| Documentation | Basic | ✅ Comprehensive guides |

---

## 🐛 Troubleshooting

### "Provider not found"
- Check `~/.terraformrc` has correct path
- Run `echo $(go env GOPATH)/bin` and use that path
- Run `make install` again

### "Dev overrides in use" warning
- This is **normal** with local development
- It's Terraform telling you it's using your local build

### Changes don't appear in Terraform
```bash
make install
rm -rf .terraform .terraform.lock.hcl
terraform init
terraform plan
```

---

## 🔗 Quick Links

- **Your Fork**: https://github.com/patrikwm/terraform-provider-coolify
- **GitHub Actions**: https://github.com/patrikwm/terraform-provider-coolify/actions
- **Create Release**: https://github.com/patrikwm/terraform-provider-coolify/releases/new
- **Upstream**: https://github.com/SierraJC/terraform-provider-coolify
- **Coolify API**: https://coolify.io/docs/api-reference/introduction

---

## 💡 Pro Tips

1. **Use dev overrides during development** - Make changes, run `make install`, changes reflect immediately
2. **Version your releases** - Follow semantic versioning (v0.1.0, v0.2.0, etc.)
3. **Test thoroughly** - Run `make test` before every commit
4. **Keep OpenAPI updated** - Run `make fetch-schema` regularly
5. **Document changes** - Update examples/ when adding features

---

## ❓ FAQ

**Q: Do I need to publish to Terraform Registry?**
A: No! Local dev or GitHub releases are fine for personal use.

**Q: Do I need GPG signing?**
A: Only if publishing to Terraform Registry. GitHub releases work without it.

**Q: Can I contribute back to upstream?**
A: Yes! Your environment resource and other enhancements could be upstreamed.

**Q: How do I update from upstream?**
A:
```bash
git remote add upstream https://github.com/SierraJC/terraform-provider-coolify.git
git fetch upstream
git merge upstream/main
```

---

**Ready to go! Pick Path A (local) or Path B (publish) above and start using your enhanced provider.** 🎉

Questions? Check [SETUP.md](SETUP.md) or [CONTRIBUTING.md](CONTRIBUTING.md).
