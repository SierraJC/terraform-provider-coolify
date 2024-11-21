<p align="center">
  <a href="https://github.com/SierraJC/terraform-provider-coolify/blob/main/LICENSE" alt="License">
    <img src="https://img.shields.io/github/license/SierraJC/terraform-provider-coolify" /></a>
  <a href="https://GitHub.com/SierraJC/terraform-provider-coolify/releases/" alt="Release">
    <img src="https://img.shields.io/github/v/release/SierraJC/terraform-provider-coolify?include_prereleases" /></a>
  <a href="https://github.com/coollabsio/coolify" alt="Coolify">
    <img src="https://img.shields.io/badge/Coolify-v4.0.0--beta.367-orange" /></a>
  <br/>
  <a href="http://golang.org" alt="Made With Go">
    <img src="https://img.shields.io/github/go-mod/go-version/SierraJC/terraform-provider-coolify" /></a>
  <a href="https://github.com/SierraJC/terraform-provider-coolify/actions/workflows/test.yml" alt="Tests">
    <img src="https://github.com/SierraJC/terraform-provider-coolify/actions/workflows/test.yml/badge.svg?branch=main" /></a>
  <a href="https://codecov.io/gh/SierraJC/terraform-provider-coolify" alt="Coverage">
    <img src="https://codecov.io/gh/SierraJC/terraform-provider-coolify/graph/badge.svg?token=63aeH0TuP2" /></a>
</p>

# Terraform Provider for [Coolify](https://coolify.io/) _v4_

Documentation: https://registry.terraform.io/providers/SierraJC/coolify/latest/docs

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

## Supported Coolify Resources

| Feature                    | Resource | Data Source |
| -------------------------- | -------- | ----------- |
| Teams                      | ⛔       | ️✔️         |
| Private Keys               | ✔️       | ✔️          |
| Servers                    | ✔️       | ️✔️         |
| - Server Resources         |          | ️✔️         |
| - Server Domains           |          | ️✔️         |
| Projects                   | ✔️       | ✔️          |
| - Project Environments     | ⛔       | ⛔          |
| Resources                  | ⛔       | ➖          |
| Databases                  | ➖       | ➖          |
| Services                   | ➖       | ⚒️          |
| - Service Environments     | ✔️       | ➖          |
| Applications               | ➖       | ⚒️          |
| - Application Environments | ✔️       | ➖          |

✔️ Supported ⚒️ Partial Support ➖ Planned ⛔ Blocked by Coolify API

The provider is currently limited by the [Coolify API](https://github.com/coollabsio/coolify/blob/main/openapi.yaml), which is still in development. As the API matures, more resources will be added to the provider.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please read the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
