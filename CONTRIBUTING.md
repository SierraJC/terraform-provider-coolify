# Contributing to terraform-provider-coolify

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- [Go](https://golang.org/doc/install) 1.24+ (check `go.mod` for exact version)
- [Terraform](https://developer.hashicorp.com/terraform/install) 1.5+ or [OpenTofu](https://opentofu.org/docs/intro/install/)
- A running [Coolify v4](https://coolify.io/) instance with API access
- Git

### Initial Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/terraform-provider-coolify.git
   cd terraform-provider-coolify
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your Coolify instance details:
   # COOLIFY_ENDPOINT=https://your-coolify.example.com/api/v1
   # COOLIFY_TOKEN=your-api-token
   ```

4. **Build the provider:**
   ```bash
   make install
   ```

## Development Workflow

### Code Generation

This provider uses code generation from OpenAPI specs. When modifying the provider:

1. **Update OpenAPI specs** (if adding/modifying Coolify API endpoints):
   ```bash
   # Fetch latest from Coolify
   make fetch-schema

   # Or manually edit tools/openapi.yml
   ```

2. **Regenerate code:**
   ```bash
   make generate
   ```

   This generates:
   - `internal/api/api_gen.go` - API client from OpenAPI spec
   - `internal/provider/generated/` - Provider schemas
   - `docs/` - Provider documentation

### Testing

```bash
# Run unit tests
make test

# Run acceptance tests (requires .env with working Coolify instance)
make testacc

# Run specific test
go test -v -run TestAccEnvironmentResource_basic ./...
```

### Local Development

1. **Build and install locally:**
   ```bash
   make install
   ```

2. **Configure Terraform to use your local build:**

   Create/edit `~/.terraformrc`:
   ```hcl
   provider_installation {
     dev_overrides {
       "registry.terraform.io/patrikwm/coolify" = "/Users/YOUR_USERNAME/go/bin"
     }
     direct {}
   }
   ```

   Replace path with your `$(go env GOPATH)/bin`.

3. **Test in a Terraform project:**
   ```hcl
   terraform {
     required_providers {
       coolify = {
         source = "patrikwm/coolify"
       }
     }
   }
   ```

   Terraform will now use your locally built version.

## Adding New Resources

### 1. Check the Coolify OpenAPI Spec

Ensure the resource exists in the [Coolify API](https://github.com/coollabsio/coolify/blob/main/openapi.yaml):

```bash
# Update local copy
make fetch-schema
# Check tools/openapi.yml
```

### 2. Add/Update OpenAPI Definitions

If endpoints are missing or incomplete in `tools/openapi.yml`, add them following the existing patterns.

### 3. Create Resource Files

```bash
# Example: coolify_foo resource
touch internal/service/foo_resource.go
touch internal/service/foo_resource_test.go
```

**Resource structure:**
```go
package service

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    // ... other imports
)

type fooResource struct {
    client *api.ClientWithResponses
}

func NewFooResource() resource.Resource {
    return &fooResource{}
}

// Implement: Metadata, Schema, Configure, Create, Read, Update, Delete, ImportState
```

### 4. Register the Resource

Edit `internal/provider/provider.go`:

```go
func (p *coolifyProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // ... existing resources
        service.NewFooResource,
    }
}
```

### 5. Write Tests

```go
package service

import (
    "testing"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFooResource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
                    resource "coolify_foo" "test" {
                        name = "test-foo"
                    }
                `,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("coolify_foo.test", "name", "test-foo"),
                ),
            },
        },
    })
}
```

### 6. Generate Documentation

```bash
make generate
# Check docs/resources/foo.md
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting: `make fmt`
- Run linters (when enabled): `make lint`
- Add descriptive comments for exported types and functions
- Keep functions focused and testable

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add coolify_environment resource` - New features
- `fix: handle nil pointer in service status` - Bug fixes
- `docs: update README with new features` - Documentation
- `test: add acceptance tests for environments` - Tests
- `refactor: simplify destination_uuid logic` - Code refactoring
- `chore: update dependencies` - Maintenance

**Breaking changes:** Add `!` after type:
```
feat!: make destination_uuid truly optional
```

## Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes:**
   - Update code
   - Add/update tests
   - Run `make generate` if needed
   - Run `make test`

3. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add my new feature"
   ```

4. **Push and create PR:**
   ```bash
   git push origin feature/my-new-feature
   ```

5. **Ensure CI passes:**
   - All tests must pass
   - Code generation must be up-to-date
   - No linter errors

## Release Process

Releases are automated via GitHub Actions:

1. **Update version** (follow [Semantic Versioning](https://semver.org/)):
   - `v0.x.y` - Breaking changes increment `x`
   - `v0.1.y` - New features increment last digit
   - `v0.1.z` - Bug fixes increment last digit

2. **Create and push tag:**
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```

3. **GitHub Action builds binaries** for all platforms and creates a draft release

4. **Review and publish** the release on GitHub

### GPG Signing (Required for Terraform Registry)

If publishing to the Terraform Registry, releases must be signed:

1. **Generate GPG key:**
   ```bash
   gpg --full-generate-key
   # Use RSA, 4096 bits, name/email matching GitHub
   ```

2. **Export keys:**
   ```bash
   gpg --armor --export YOUR_EMAIL > public.key
   gpg --armor --export-secret-keys YOUR_EMAIL > private.key
   ```

3. **Add secrets to GitHub:**
   - `GPG_PRIVATE_KEY` - Contents of `private.key`
   - `PASSPHRASE` - GPG key passphrase
   - `GPG_FINGERPRINT` - Your key fingerprint

See [Terraform Registry docs](https://developer.hashicorp.com/terraform/registry/providers/publishing#preparing-and-adding-a-signing-key) for details.

## Architecture Overview

```
terraform-provider-coolify/
├── internal/
│   ├── api/              # Auto-generated API client
│   │   └── api_gen.go    # From tools/openapi.yml
│   ├── provider/         # Provider implementation
│   │   ├── provider.go   # Main provider logic
│   │   └── generated/    # Auto-generated schemas
│   └── service/          # Resource/data source implementations
│       ├── *_resource.go      # Resource CRUD
│       ├── *_data_source.go   # Data sources
│       └── *_test.go          # Tests
├── tools/
│   ├── openapi.yml       # Coolify API spec (base)
│   ├── overlay.yml       # Modifications to upstream spec
│   └── tfplugingen-*.yml # Code generation config
├── docs/                 # Auto-generated documentation
└── examples/             # Usage examples
```

## Getting Help

- Check existing [GitHub Issues](https://github.com/patrikwm/terraform-provider-coolify/issues)
- Review the [Coolify API documentation](https://coolify.io/docs/api-reference/introduction)
- Look at [upstream provider](https://github.com/SierraJC/terraform-provider-coolify) for context

## Code of Conduct

Be respectful and constructive. This is a community project.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

