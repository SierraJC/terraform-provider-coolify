# Coolify API Coverage — Terraform Provider

> Generated: 2026-03-01 | Based on Coolify OpenAPI v0.1 (107 endpoints)

## Coverage Summary

| Category | Total Endpoints | Used | Unused | Coverage |
|----------|----------------|------|--------|----------|
| Applications | 19 | 5 | 14 | 26% |
| Cloud Tokens | 6 | 0 | 6 | 0% |
| Databases | 21 | 7 | 14 | 33% |
| Deployments | 5 | 0 | 5 | 0% |
| GitHub Apps | 6 | 0 | 6 | 0% |
| Hetzner | 5 | 0 | 5 | 0% |
| General | 4 | 1 | 3 | 25% |
| Projects | 9 | 8 | 1 | 89% |
| Resources | 1 | 0 | 1 | 0% |
| Private Keys | 5 | 5 | 0 | **100%** |
| Servers | 8 | 6 | 2 | 75% |
| Services | 13 | 11 | 2 | 85% |
| Teams | 5 | 4 | 1 | 80% |
| **Total** | **107** | **47** | **60** | **44%** |

---

## Currently Implemented (9 Resources + 13 Data Sources)

### Resources

| Resource | CRUD | API Endpoints |
|----------|------|---------------|
| `coolify_server` | Create, Read, Update, Delete, Import | `POST /servers`, `GET /servers/{uuid}`, `PATCH /servers/{uuid}`, `DELETE /servers/{uuid}` |
| `coolify_private_key` | Create, Read, Update, Delete, Import | `POST /security/keys`, `GET /security/keys/{uuid}`, `PATCH /security/keys/{uuid}`, `DELETE /security/keys/{uuid}` |
| `coolify_project` | Create, Read, Update, Delete, Import | `POST /projects`, `GET /projects/{uuid}`, `PATCH /projects/{uuid}`, `DELETE /projects/{uuid}` |
| `coolify_environment` | Create, Read, Delete, Import | `POST /projects/{uuid}/environments`, `GET /projects/{uuid}/{env}`, `DELETE /projects/{uuid}/{env}` |
| `coolify_service` | Create, Read, Update, Delete, Import | `POST /services`, `GET /services/{uuid}`, `PATCH /services/{uuid}`, `DELETE /services/{uuid}`, `GET /services/{uuid}/restart` |
| `coolify_service_envs` | Create, Read, Update, Delete, Import | `POST /services/{uuid}/envs`, `GET /services/{uuid}/envs`, `PATCH /services/{uuid}/envs/bulk`, `DELETE /services/{uuid}/envs/{env_uuid}` |
| `coolify_application_envs` | Create, Read, Update, Delete, Import | `POST /applications/{uuid}/envs`, `GET /applications/{uuid}/envs`, `PATCH /applications/{uuid}/envs/bulk`, `DELETE /applications/{uuid}/envs/{env_uuid}` |
| `coolify_postgresql_database` | Create, Read, Update, Delete, Import | `POST /databases/postgresql`, `GET /databases/{uuid}`, `PATCH /databases/{uuid}`, `DELETE /databases/{uuid}`, `GET /databases/{uuid}/restart` |
| `coolify_mysql_database` | Create, Read, Update, Delete, Import | `POST /databases/mysql`, `GET /databases/{uuid}`, `PATCH /databases/{uuid}`, `DELETE /databases/{uuid}`, `GET /databases/{uuid}/restart` |

### Data Sources

| Data Source | API Endpoints |
|-------------|---------------|
| `coolify_server` | `GET /servers/{uuid}` |
| `coolify_servers` | `GET /servers` |
| `coolify_server_resources` | `GET /servers/{uuid}/resources` |
| `coolify_server_domains` | `GET /servers/{uuid}/domains` |
| `coolify_private_key` | `GET /security/keys/{uuid}` |
| `coolify_private_keys` | `GET /security/keys` |
| `coolify_project` | `GET /projects/{uuid}` |
| `coolify_projects` | `GET /projects` |
| `coolify_application` | `GET /applications/{uuid}` |
| `coolify_applications` | `GET /applications` |
| `coolify_service` | `GET /services/{uuid}` |
| `coolify_team` | `GET /teams/{id}`, `GET /teams/current`, `GET /teams/{id}/members` |
| `coolify_teams` | `GET /teams`, `GET /teams/{id}/members` |

---

## Missing — Grouped by Implementation Priority

### Phase 1 — Database Resources (Quick Wins)

Clone the existing `coolify_mysql_database` pattern. Each uses `POST /databases/{type}` for creation plus the shared `GET/PATCH/DELETE /databases/{uuid}` endpoints.

| New Resource | Create Endpoint | Effort |
|-------------|-----------------|--------|
| `coolify_mariadb_database` | `POST /databases/mariadb` | Small |
| `coolify_mongodb_database` | `POST /databases/mongodb` | Small |
| `coolify_redis_database` | `POST /databases/redis` | Small |
| `coolify_clickhouse_database` | `POST /databases/clickhouse` | Small |
| `coolify_keydb_database` | `POST /databases/keydb` | Small |
| `coolify_dragonfly_database` | `POST /databases/dragonfly` | Small |

**Endpoints covered:** 6 new (create) + 18 shared (CRUD × 6) = brings Databases from 33% → **100%** (excluding backups)

---

### Phase 2 — Missing Data Sources

Simple read-only resources with high user value.

| New Data Source | API Endpoint | Effort |
|----------------|-------------|--------|
| `coolify_databases` | `GET /databases` | Small |
| `coolify_services` | `GET /services` | Small |
| `coolify_environment` | `GET /projects/{uuid}/{env}` | Small |
| `coolify_environments` | `GET /projects/{uuid}/environments` | Small |
| `coolify_resources` | `GET /resources` | Small |

**Endpoints covered:** 5 new

---

### Phase 3 — Application Resource (Biggest Impact)

The most complex but most valuable missing piece. The API supports 5 different creation methods.

| Feature | API Endpoints | Notes |
|---------|---------------|-------|
| Create (Public Git) | `POST /applications/public` | Most common use case |
| Create (Private - GH App) | `POST /applications/private-github-app` | Requires github_app_id |
| Create (Private - Deploy Key) | `POST /applications/private-deploy-key` | Requires private_key_uuid |
| Create (Dockerfile) | `POST /applications/dockerfile` | No git required |
| Create (Docker Image) | `POST /applications/dockerimage` | No git required |
| Create (Docker Compose) | `POST /applications/dockercompose` | **Deprecated** — use services |
| Read | `GET /applications/{uuid}` | Already in data source |
| Update | `PATCH /applications/{uuid}` | ~100 updatable fields |
| Delete | `DELETE /applications/{uuid}` | With optional cleanup params |
| Start | `GET /applications/{uuid}/start` | Operational |
| Stop | `GET /applications/{uuid}/stop` | Operational |
| Restart | `GET /applications/{uuid}/restart` | Operational |
| Logs | `GET /applications/{uuid}/logs` | Read-only, data source candidate |

**Effort:** Large — Multiple creation methods, extensive schema, lifecycle management
**Endpoints covered:** Up to 12 (excluding deprecated dockercompose)

---

### Phase 4 — Database Backups

| New Resource | API Endpoints | Effort |
|-------------|---------------|--------|
| `coolify_database_backup` | `POST /databases/{uuid}/backups` (create) | Medium |
| | `GET /databases/{uuid}/backups` (list) | |
| | `PATCH /databases/{uuid}/backups/{id}` (update) | |
| | `DELETE /databases/{uuid}/backups/{id}` (delete) | |
| `coolify_database_backup_executions` (data source) | `GET /databases/{uuid}/backups/{id}/executions` (list) | Small |
| | `DELETE /databases/{uuid}/backups/{id}/executions/{exec_id}` (delete) | |

**Endpoints covered:** 6

---

### Phase 5 — GitHub Apps

| New Resource | API Endpoints | Effort |
|-------------|---------------|--------|
| `coolify_github_app` | `POST /github-apps` (create) | Medium |
| | `GET /github-apps` (list — data source) | |
| | `PATCH /github-apps/{id}` (update) | |
| | `DELETE /github-apps/{id}` (delete) | |
| `coolify_github_app_repositories` (data source) | `GET /github-apps/{id}/repositories` | Small |
| `coolify_github_app_branches` (data source) | `GET /github-apps/{id}/repositories/{owner}/{repo}/branches` | Small |

**Endpoints covered:** 6

---

### Phase 6 — Cloud Tokens

| New Resource | API Endpoints | Effort |
|-------------|---------------|--------|
| `coolify_cloud_token` | `POST /cloud-tokens` (create) | Medium |
| | `GET /cloud-tokens/{uuid}` (read) | |
| | `PATCH /cloud-tokens/{uuid}` (update) | |
| | `DELETE /cloud-tokens/{uuid}` (delete) | |
| `coolify_cloud_tokens` (data source) | `GET /cloud-tokens` (list) | Small |
| Validate action | `POST /cloud-tokens/{uuid}/validate` | Optional |

**Endpoints covered:** 6

---

### Phase 7 — Hetzner Integration

| New Resource/Data Source | API Endpoints | Effort |
|-------------------------|---------------|--------|
| `coolify_hetzner_server` | `POST /servers/hetzner` (create) | Medium |
| `coolify_hetzner_locations` (data source) | `GET /hetzner/locations` | Small |
| `coolify_hetzner_server_types` (data source) | `GET /hetzner/server-types` | Small |
| `coolify_hetzner_images` (data source) | `GET /hetzner/images` | Small |
| `coolify_hetzner_ssh_keys` (data source) | `GET /hetzner/ssh-keys` | Small |

**Endpoints covered:** 5

---

### Phase 8 — Operational / Utility

These are less typical for Terraform (actions rather than state) but could be implemented as resources or provider-level features.

| Feature | API Endpoints | Implementation |
|---------|---------------|----------------|
| Start/Stop/Restart applications | `GET /applications/{uuid}/start\|stop\|restart` | Attribute on application resource |
| Start/Stop/Restart databases | `GET /databases/{uuid}/start\|stop\|restart` | Attribute on database resources |
| Start/Stop services | `GET /services/{uuid}/start\|stop` | Attribute on service resource (restart exists) |
| Deploy by tag/UUID | `GET /deploy` | Operational — `null_resource` with provisioner |
| Cancel deployment | `POST /deployments/{uuid}/cancel` | Operational |
| Deployment data | `GET /deployments`, `GET /deployments/{uuid}`, `GET /deployments/applications/{uuid}` | Data source |
| Server validate | `GET /servers/{uuid}/validate` | Attribute on server resource |
| API enable/disable | `GET /enable`, `GET /disable` | Provider-level, unlikely needed |
| Health check | `GET /health` | Provider-level, unlikely needed |
| Version | `GET /version` | Already used by provider internally |

**Endpoints covered:** 15+ (mostly operational)

---

## Full Endpoint Reference

### ✅ = Implemented | ❌ = Not Implemented | ⚠️ = Deprecated

#### Applications (19 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/applications` | `list-applications` |
| ❌ | POST | `/applications/public` | `create-public-application` |
| ❌ | POST | `/applications/private-github-app` | `create-private-github-app-application` |
| ❌ | POST | `/applications/private-deploy-key` | `create-private-deploy-key-application` |
| ❌ | POST | `/applications/dockerfile` | `create-dockerfile-application` |
| ❌ | POST | `/applications/dockerimage` | `create-dockerimage-application` |
| ⚠️ | POST | `/applications/dockercompose` | `create-dockercompose-application` |
| ✅ | GET | `/applications/{uuid}` | `get-application-by-uuid` |
| ❌ | DELETE | `/applications/{uuid}` | `delete-application-by-uuid` |
| ❌ | PATCH | `/applications/{uuid}` | `update-application-by-uuid` |
| ❌ | GET | `/applications/{uuid}/logs` | `get-application-logs-by-uuid` |
| ✅ | GET | `/applications/{uuid}/envs` | `list-envs-by-application-uuid` |
| ✅ | POST | `/applications/{uuid}/envs` | `create-env-by-application-uuid` |
| ❌ | PATCH | `/applications/{uuid}/envs` | `update-env-by-application-uuid` |
| ✅ | PATCH | `/applications/{uuid}/envs/bulk` | `update-envs-by-application-uuid` |
| ✅ | DELETE | `/applications/{uuid}/envs/{env_uuid}` | `delete-env-by-application-uuid` |
| ❌ | GET | `/applications/{uuid}/start` | `start-application-by-uuid` |
| ❌ | GET | `/applications/{uuid}/stop` | `stop-application-by-uuid` |
| ❌ | GET | `/applications/{uuid}/restart` | `restart-application-by-uuid` |

#### Cloud Tokens (6 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/cloud-tokens` | `list-cloud-tokens` |
| ❌ | POST | `/cloud-tokens` | `create-cloud-token` |
| ❌ | GET | `/cloud-tokens/{uuid}` | `get-cloud-token-by-uuid` |
| ❌ | DELETE | `/cloud-tokens/{uuid}` | `delete-cloud-token-by-uuid` |
| ❌ | PATCH | `/cloud-tokens/{uuid}` | `update-cloud-token-by-uuid` |
| ❌ | POST | `/cloud-tokens/{uuid}/validate` | `validate-cloud-token-by-uuid` |

#### Databases (21 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/databases` | `list-databases` |
| ✅ | GET | `/databases/{uuid}` | `get-database-by-uuid` |
| ✅ | DELETE | `/databases/{uuid}` | `delete-database-by-uuid` |
| ✅ | PATCH | `/databases/{uuid}` | `update-database-by-uuid` |
| ✅ | POST | `/databases/postgresql` | `create-database-postgresql` |
| ❌ | POST | `/databases/clickhouse` | `create-database-clickhouse` |
| ❌ | POST | `/databases/dragonfly` | `create-database-dragonfly` |
| ❌ | POST | `/databases/redis` | `create-database-redis` |
| ❌ | POST | `/databases/keydb` | `create-database-keydb` |
| ❌ | POST | `/databases/mariadb` | `create-database-mariadb` |
| ✅ | POST | `/databases/mysql` | `create-database-mysql` |
| ❌ | POST | `/databases/mongodb` | `create-database-mongodb` |
| ❌ | GET | `/databases/{uuid}/backups` | `get-database-backups-by-uuid` |
| ❌ | POST | `/databases/{uuid}/backups` | `create-database-backup` |
| ❌ | DELETE | `/databases/{uuid}/backups/{id}` | `delete-backup-configuration-by-uuid` |
| ❌ | PATCH | `/databases/{uuid}/backups/{id}` | `update-database-backup` |
| ❌ | GET | `/databases/{uuid}/backups/{id}/executions` | `list-backup-executions` |
| ❌ | DELETE | `/databases/{uuid}/backups/{id}/executions/{exec_id}` | `delete-backup-execution-by-uuid` |
| ❌ | GET | `/databases/{uuid}/start` | `start-database-by-uuid` |
| ❌ | GET | `/databases/{uuid}/stop` | `stop-database-by-uuid` |
| ✅ | GET | `/databases/{uuid}/restart` | `restart-database-by-uuid` |

#### Deployments (5 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/deployments` | `list-deployments` |
| ❌ | GET | `/deployments/{uuid}` | `get-deployment-by-uuid` |
| ❌ | POST | `/deployments/{uuid}/cancel` | `cancel-deployment-by-uuid` |
| ❌ | GET | `/deploy` | `deploy-by-tag-or-uuid` |
| ❌ | GET | `/deployments/applications/{uuid}` | `list-deployments-by-app-uuid` |

#### GitHub Apps (6 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/github-apps` | `list-github-apps` |
| ❌ | POST | `/github-apps` | `create-github-app` |
| ❌ | GET | `/github-apps/{id}/repositories` | `load-repositories` |
| ❌ | GET | `/github-apps/{id}/repositories/{owner}/{repo}/branches` | `load-branches` |
| ❌ | DELETE | `/github-apps/{id}` | `deleteGithubApp` |
| ❌ | PATCH | `/github-apps/{id}` | `updateGithubApp` |

#### Hetzner (5 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/hetzner/locations` | `get-hetzner-locations` |
| ❌ | GET | `/hetzner/server-types` | `get-hetzner-server-types` |
| ❌ | GET | `/hetzner/images` | `get-hetzner-images` |
| ❌ | GET | `/hetzner/ssh-keys` | `get-hetzner-ssh-keys` |
| ❌ | POST | `/servers/hetzner` | `create-hetzner-server` |

#### General (4 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/version` | `version` |
| ❌ | GET | `/enable` | `enable-api` |
| ❌ | GET | `/disable` | `disable-api` |
| ❌ | GET | `/health` | `healthcheck` |

#### Projects (9 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/projects` | `list-projects` |
| ✅ | POST | `/projects` | `create-project` |
| ✅ | GET | `/projects/{uuid}` | `get-project-by-uuid` |
| ✅ | DELETE | `/projects/{uuid}` | `delete-project-by-uuid` |
| ✅ | PATCH | `/projects/{uuid}` | `update-project-by-uuid` |
| ✅ | GET | `/projects/{uuid}/{env}` | `get-environment-by-name-or-uuid` |
| ✅ | GET | `/projects/{uuid}/environments` | `get-environments` |
| ✅ | POST | `/projects/{uuid}/environments` | `create-environment` |
| ✅ | DELETE | `/projects/{uuid}/environments/{env}` | `delete-environment` |

#### Resources (1 endpoint)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/resources` | `list-resources` |

#### Private Keys (5 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/security/keys` | `list-private-keys` |
| ✅ | POST | `/security/keys` | `create-private-key` |
| ✅ | PATCH | `/security/keys` | `update-private-key` |
| ✅ | GET | `/security/keys/{uuid}` | `get-private-key-by-uuid` |
| ✅ | DELETE | `/security/keys/{uuid}` | `delete-private-key-by-uuid` |

#### Servers (8 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/servers` | `list-servers` |
| ✅ | POST | `/servers` | `create-server` |
| ✅ | GET | `/servers/{uuid}` | `get-server-by-uuid` |
| ✅ | DELETE | `/servers/{uuid}` | `delete-server-by-uuid` |
| ✅ | PATCH | `/servers/{uuid}` | `update-server-by-uuid` |
| ✅ | GET | `/servers/{uuid}/resources` | `get-resources-by-server-uuid` |
| ✅ | GET | `/servers/{uuid}/domains` | `get-domains-by-server-uuid` |
| ❌ | GET | `/servers/{uuid}/validate` | `validate-server-by-uuid` |

#### Services (13 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ❌ | GET | `/services` | `list-services` |
| ✅ | POST | `/services` | `create-service` |
| ✅ | GET | `/services/{uuid}` | `get-service-by-uuid` |
| ✅ | DELETE | `/services/{uuid}` | `delete-service-by-uuid` |
| ✅ | PATCH | `/services/{uuid}` | `update-service-by-uuid` |
| ✅ | GET | `/services/{uuid}/envs` | `list-envs-by-service-uuid` |
| ✅ | POST | `/services/{uuid}/envs` | `create-env-by-service-uuid` |
| ❌ | PATCH | `/services/{uuid}/envs` | `update-env-by-service-uuid` |
| ✅ | PATCH | `/services/{uuid}/envs/bulk` | `update-envs-by-service-uuid` |
| ✅ | DELETE | `/services/{uuid}/envs/{env_uuid}` | `delete-env-by-service-uuid` |
| ❌ | GET | `/services/{uuid}/start` | `start-service-by-uuid` |
| ❌ | GET | `/services/{uuid}/stop` | `stop-service-by-uuid` |
| ✅ | GET | `/services/{uuid}/restart` | `restart-service-by-uuid` |

#### Teams (5 endpoints)

| Status | Method | Path | operationId |
|--------|--------|------|-------------|
| ✅ | GET | `/teams` | `list-teams` |
| ✅ | GET | `/teams/{id}` | `get-team-by-id` |
| ✅ | GET | `/teams/{id}/members` | `get-members-by-team-id` |
| ✅ | GET | `/teams/current` | `get-current-team` |
| ✅ | GET | `/teams/current/members` | `get-current-team-members` |
