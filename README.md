<p align="center">
  <a href="https://github.com/patrikwm/terraform-provider-coolify/blob/main/LICENSE" alt="License">
    <img src="https://img.shields.io/github/license/patrikwm/terraform-provider-coolify" /></a>
  <a href="https://GitHub.com/patrikwm/terraform-provider-coolify/releases/" alt="Release">
    <img src="https://img.shields.io/github/v/release/patrikwm/terraform-provider-coolify?include_prereleases" /></a>
  <a href="https://github.com/coollabsio/coolify" alt="Coolify">
    <img src="https://img.shields.io/badge/Coolify-v4.x-orange" /></a>
  <br/>
  <a href="http://golang.org" alt="Made With Go">
    <img src="https://img.shields.io/github/go-mod/go-version/patrikwm/terraform-provider-coolify" /></a>
  <a href="https://github.com/patrikwm/terraform-provider-coolify/actions/workflows/test.yml" alt="Tests">
    <img src="https://github.com/patrikwm/terraform-provider-coolify/actions/workflows/test.yml/badge.svg?branch=main" /></a>
</p>

# Terraform Provider for [Coolify](https://coolify.io/) _v4_

_This project is a community-driven fork with enhanced features. Not affiliated with or an official product of Coolify._

**Upstream:** Based on [SierraJC/terraform-provider-coolify](https://github.com/SierraJC/terraform-provider-coolify) with additional enhancements.

Documentation: See [docs/](docs/) directory or https://registry.terraform.io/providers/SierraJC/coolify/latest/docs (upstream)

The Coolify provider enables Terraform to manage [Coolify](https://coolify.io/) _v4 (beta)_ resources.
See the [examples](examples/) directory for usage examples.

This project follows [Semantic Versioning](https://semver.org/). As the current version is 0.x.x, the API should be considered unstable and subject to breaking changes.

## Prerequisites

Before you begin using the Coolify Terraform Provider, ensure you have completed the following steps:

1. Install Terraform by following the official [HashiCorp documentation](https://developer.hashicorp.com/terraform/install).
1. Create a new API token with _Root Access_ in the Coolify dashboard. See the [Coolify API documentation](https://coolify.io/docs/api-reference/authorization#generate)
1. Set the `COOLIFY_TOKEN` environment variable to your API token. For example, add the following line to your `.bashrc` file:
   ```bash
   export COOLIFY_TOKEN="Your API token"
   ```

**For detailed setup instructions, see [SETUP.md](SETUP.md).**

## ✨ Enhanced Features in This Fork

- ✅ **`coolify_environment` resource** - Full CRUD for project environments (upstream blocked)
- ✅ **Optional `destination_uuid`** - Services auto-select destinations when omitted
- ✅ **Computed status attributes** - Services expose `status` and `server_status` from API
- 📦 **Updated OpenAPI schemas** - Environment and Service models match actual API responses

## Supported Coolify Resources

| Feature                    | Resource | Data Source | Notes                                |
| -------------------------- | -------- | ----------- | ------------------------------------ |
| Teams                      | ⛔       | ️✔️         |                                      |
| Private Keys               | ✔️       | ✔️          |                                      |
| Servers                    | ✔️       | ️✔️         |                                      |
| - Server Resources         |          | ️✔️         |                                      |
| - Server Domains           |          | ️✔️         |                                      |
| Destinations               | ⛔       | ⛔          |                                      |
| Projects                   | ✔️       | ✔️          |                                      |
| - Project Environments     | ✔️       | ⛔          | **✨ New in this fork**              |
| Resources                  | ⛔       | ⛔          |                                      |
| Databases                  | ⚒️       | ➖          | PostgreSQL & MySQL only              |
| Services                   | ✔️       | ⚒️          | **✨ Enhanced with status fields**   |
| - Service Environments     | ✔️       | ➖          |                                      |
| Applications               | ⚒️       | ✔️          |                                      |
| - Application Environments | ✔️       | ➖          |                                      |

✔️ Supported ⚒️ Partial Support ➖ Planned ⛔ Blocked by Coolify API

The provider is currently limited by the [Coolify API](https://github.com/coollabsio/coolify/blob/main/openapi.yaml), which is still in development. As the API matures, more resources will be added to the provider.

## Using This Provider

### Option 1: Local Development (Recommended for Testing)

1. **Build the provider:**
   ```bash
   make install
   ```

2. **Configure Terraform to use the local build:**

   Create/edit `~/.terraformrc` (or `%APPDATA%/terraform.rc` on Windows):
   ```hcl
   provider_installation {
     dev_overrides {
       "registry.terraform.io/patrikwm/coolify" = "/Users/patrik/go/bin"
     }
     direct {}
   }
   ```

   Replace the path with your actual `$(go env GOPATH)/bin` directory.

3. **Use in your Terraform configuration:**
   ```hcl
   terraform {
     required_providers {
       coolify = {
         source = "patrikwm/coolify"
       }
     }
   }

   provider "coolify" {
     endpoint = "https://your-coolify.example.com/api/v1"
     # token can be set via COOLIFY_TOKEN env var
   }
   ```

### Option 2: Use Published GitHub Release

1. **Create a GitHub release** with a version tag (e.g., `v0.1.0`)
2. **The release workflow** will automatically build binaries for all platforms
3. **Configure Terraform:**
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

### Option 3: Publish to Terraform Registry (Optional)

Follow the [Terraform Registry publishing guide](https://developer.hashicorp.com/terraform/registry/providers/publishing) to make the provider available via `registry.terraform.io`.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development instructions.

**Quick start:**
```bash
make generate  # Generate code from OpenAPI spec
make test      # Run unit tests
make testacc   # Run acceptance tests (requires .env with COOLIFY_ENDPOINT and COOLIFY_TOKEN)
make install   # Build and install to GOPATH
```

## Contributing

Contributions are welcome! If you would like to contribute to this project, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
