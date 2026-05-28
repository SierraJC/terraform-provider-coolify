# Changelog

All notable changes to this fork are documented here. This fork is based on
`SierraJC/terraform-provider-coolify` with `coolify_application` resource added
from PR #87 plus subsequent fixes.

## v1.1.0 (2026-05-28)

### Added
- `CoolifyLabelsSemanticEqual` plan modifier on `custom_labels`. Recognizes
  Coolify v4's server-side label normalization (base64↔plaintext re-encoding +
  automatic `tls.certresolver=letsencrypt` injection on Traefik routers) as
  semantic no-op. Prevents drift loops; enables safe Tofu adoption of existing
  Coolify applications with file-based TLS setups (e.g., CF Origin Cert via
  Traefik dynamic config).
- Unit test coverage for normalization algorithm (12 table-driven cases +
  isolated helper tests).

### Notes
- `ToAPIUpdate` continues to omit `custom_labels` (carried over from v1.0.3) —
  defensive, avoids unnecessary server-side mutation when only other fields
  change.
- Backward-compatible: no schema breakage. Downgrade NOT supported (state
  written by v1.1.0 remains compatible, but v1.0.x will resume showing drift).

## v1.0.3 (2026-05-28)

### Fixed
- Omit `custom_labels` from `UpdateApplicationByUuidJSONRequestBody`. Coolify
  v4 mutates labels on update; sending them explicitly previously overwrote
  user intent. (Mitigation attempt — Coolify still mutates server-side even
  with payload omission; full fix arrives in v1.1.0.)

## v1.0.2 (2026-05-28)

### Fixed
- Omit `destination_uuid` from update payload (create-only on Coolify API).
- Remove `Computed: true` flag from `destination_uuid` schema (API does not
  return it on read, breaks consistency).

## v1.0.1 (2026-05-28)

### Fixed
- Omit `project_uuid`, `server_uuid`, `environment_name` from update payload.
  These fields are create-only on Coolify v4 and return HTTP 422 when sent in
  an update body.

## v1.0.0 (2026-05-28)

### Added
- `coolify_application` resource (CRUD + import + 6 source types) from
  PR #87 by FaureAlexis against upstream SierraJC/terraform-provider-coolify.
- Fork distribution under registry.wapps.co with GPG-signed releases.
