# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Environment resource implementation with full CRUD support
- Computed `status` and `server_status` attributes on service resources
- Optional `destination_uuid` on service resources (auto-selects when omitted)
- Comprehensive setup documentation (SETUP.md, QUICKSTART.md, FORK_SETUP.md)
- Enhanced CONTRIBUTING.md with detailed developer workflow

### Changed
- Made `destination_uuid` truly optional on services (removes empty string default)
- Updated OpenAPI spec to include Environment.uuid field
- Updated OpenAPI spec to include Service.status and Service.server_status fields
- Updated README.md for fork-specific features and usage

### Fixed
- Environment resource import format documented correctly: `project_uuid/environment_name_or_uuid`
- Test.yml workflow now handles missing codecov token gracefully
- GoReleaser configuration updated for patrikwm repository

## [0.1.0] - 2026-03-01

### Added
- Initial fork from SierraJC/terraform-provider-coolify
- Fork-specific enhancements and documentation

---

## Version History

This fork is based on [SierraJC/terraform-provider-coolify](https://github.com/SierraJC/terraform-provider-coolify).

For upstream changes prior to this fork, see the [upstream repository](https://github.com/SierraJC/terraform-provider-coolify/releases).

---

## Changelog Guidelines

When adding entries:

- **Added** for new features
- **Changed** for changes in existing functionality
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** in case of vulnerabilities

Include issue/PR numbers when applicable: `(#123)`

Example entry:
```markdown
### Added
- New `coolify_database` resource for managing databases (#42)

### Fixed
- Service resource now correctly handles optional destination_uuid (#43)
```
