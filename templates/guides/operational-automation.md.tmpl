---
page_title: "Operational Automation Guide"
description: |-
  Guide to using automated waiters and deployment features in the Coolify provider
---

# Operational Automation Guide

This guide explains how to use the operational automation features in this fork to eliminate manual steps in your infrastructure deployment workflow.

## Overview

Traditional infrastructure deployment with Coolify requires manual intervention:
1. Create a server → **manually click "Validate"** in UI
2. Deploy a service → **manually verify health** in UI
3. Update environment variables → **manually redeploy service** in UI

With these automation features, `terraform apply` handles everything hands-off.

## Server Validation Waiter

### Problem Solved

After creating a server in Coolify, you must manually:
- Navigate to the server in the UI
- Click the "Validate" button
- Wait for validation to complete
- Check if server is reachable and usable

### Solution

Use `wait_for_validation` to automate this:

```terraform
resource "coolify_server" "homelab" {
  name             = "Homelab Server"
  ip               = "192.168.1.100"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.main.uuid
  instant_validate = true
  
  # Automate validation - no manual UI clicks needed
  wait_for_validation = true
  
  # Optional: customize timeout (default: 3 minutes)
  timeouts {
    create = "5m"
  }
}
```

### How It Works

When `wait_for_validation = true`:
1. Server is created via API
2. Provider polls `GET /api/v1/servers/{uuid}` every 5 seconds
3. Checks `settings.is_reachable && settings.is_usable`
4. If validation fails (errors in `validation_logs`), immediate failure
5. If timeout expires before validation completes, returns timeout error
6. Only succeeds when server is fully validated

### When To Use

- **Always use for production servers** - ensures server is ready before deploying resources
- **Use with instant_validate = true** - triggers validation immediately
- **Increase timeout for slow networks** - some servers take longer to validate

### When Not To Use

- **Development/testing** - set `wait_for_validation = false` for faster iteration
- **Pre-validated servers** - if you know server is already configured

## Service Deployment Waiter

### Problem Solved

After deploying a service, you must manually:
- Navigate to the service in the UI
- Wait for containers to start
- Verify healthchecks are passing
- Confirm service is actually serving traffic

### Solution

Use `wait_for_deployment` to automate health verification:

```terraform
resource "coolify_service" "api" {
  name        = "Production API"
  server_uuid = coolify_server.homelab.uuid
  project_uuid = coolify_project.main.uuid
  environment_name = "production"
  
  instant_deploy = true
  
  # Wait for service to be healthy before completing
  wait_for_deployment = true
  
  # Customize timeouts based on your service
  timeouts {
    create = "15m"  # Initial deployment (image pull + startup)
    update = "10m"  # Redeployments (faster, using cached images)
  }
  
  compose = <<EOF
services:
  api:
    image: "myorg/api:v1.0.0"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s
EOF
}
```

### How It Works

When `wait_for_deployment = true`:
1. Service is created/updated and deployed
2. Provider polls `GET /api/v1/services/{uuid}` every 10 seconds
3. Parses `status` field (format: "status:health" like "running:healthy")
4. Success condition: `statusStr == "running" && health == "healthy"`
5. Failure conditions: `statusStr == "exited"` or `statusStr == "error"`
6. Continues polling until success, failure, or timeout

### Status Format

Coolify service status uses `status:health` format:

- ✅ **"running:healthy"** - Service is working (success)
- ⚠️ **"running:unhealthy"** - Running but failing healthchecks (keep polling)
- ⏳ **"deploying:healthy"** - Deployment in progress (keep polling)
- ❌ **"exited:unhealthy"** - Container stopped (immediate failure)
- ❌ **"error:unhealthy"** - Deployment failed (immediate failure)

### Timeout Recommendations

**Create timeout** (initial deployment):
- Simple services (e.g., nginx): 5-10 minutes
- Application services: 10-15 minutes  
- Database services: 15-20 minutes
- Large Docker images: Add 5-10 minutes for pull time

**Update timeout** (redeployment):
- Usually shorter than create (images cached)
- Default 10 minutes is sufficient for most cases
- Increase if service has long startup time

### Best Practices

1. **Always define healthchecks** in your Docker Compose:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost/health"]
  interval: 10s
  timeout: 5s
  retries: 3
  start_period: 30s  # Give app time to start before checking
```

2. **Set appropriate timeouts** based on your service:
```terraform
timeouts {
  create = "15m"  # Image pull + startup
  update = "10m"  # Just restart
}
```

3. **Use with instant_deploy**:
```terraform
instant_deploy = true
wait_for_deployment = true
```

4. **Chain dependent resources**:
```terraform
resource "coolify_service" "db" {
  wait_for_deployment = true
  # ...
}

resource "coolify_service" "api" {
  # API depends on DB being healthy
  server_uuid = coolify_server.homelab.uuid
  
  compose = <<EOF
services:
  api:
    depends_on:
      - db
EOF
  
  wait_for_deployment = true
}

resource "coolify_service_envs" "api_config" {
  # Only configure API after it's deployed
  uuid = coolify_service.api.uuid
  # ...
}
```

## Redeploy on Change

### Problem Solved

After updating environment variables, you must manually:
- Navigate to the service in the UI
- Click "Redeploy" or "Restart"
- Wait for restart to complete
- Verify service picked up new configuration

### Solution

Use `redeploy_on_change` to automate service restarts:

```terraform
resource "coolify_service_envs" "api_config" {
  uuid = coolify_service.api.uuid
  
  # Automatically restart service when any env var changes
  redeploy_on_change = true
  
  env {
    key   = "DATABASE_URL"
    value = "postgresql://user:pass@${coolify_database.db.internal_db_url}/mydb"
  }
  
  env {
    key   = "API_KEY"
    value = var.api_key
  }
  
  env {
    key   = "FEATURE_FLAGS"
    value = jsonencode({
      new_ui = true
      beta_api = false
    })
  }
}
```

### How It Works

When `redeploy_on_change = true`:
1. Environment variables are updated via API
2. Provider calls `POST /api/v1/services/{uuid}/restart`
3. If restart fails, provider issues a **warning** (not error)
4. Environment variables are still updated successfully
5. Service will use new values on next restart (manual or automatic)

### Behavior Notes

- **Warning on failure**: Redeployment failure only logs a warning, doesn't fail the apply
- **Reason**: Environment variables are the source of truth; restart is best-effort
- **Manual fallback**: If automatic restart fails, manually restart in UI

### When To Use

**Use `redeploy_on_change = true` when**:
- Environment variables affect runtime behavior (most apps)
- Service needs restart to reload configuration
- You want fully automated deployments

**Common use cases**:
- Configuration changes (API keys, URLs, feature flags)
- Database connection strings
- Service dependencies
- Runtime settings

### When Not To Use

**Use `redeploy_on_change = false` (default) when**:
- Environment variables only needed at build time
- Service hot-reloads configuration without restart
- You want to batch multiple changes before redeploying
- Testing configuration before restart

### Example: Batch Updates Without Restart

```terraform
resource "coolify_service_envs" "batch_config" {
  uuid = coolify_service.api.uuid
  
  # Don't restart yet - making multiple changes
  redeploy_on_change = false
  
  env {
    key   = "CONFIG_1"
    value = "value1"
  }
  
  env {
    key   = "CONFIG_2"
    value = "value2"
  }
  
  # To apply changes, manually trigger redeploy in UI
  # or update with redeploy_on_change = true later
}
```

## Complete Example: Fully Automated Deployment

Here's a complete example showing all automation features together:

```terraform
# Generate SSH key
resource "tls_private_key" "homelab" {
  algorithm = "ED25519"
}

# Register key with Coolify
resource "coolify_private_key" "homelab" {
  name        = "Homelab Key"
  description = "Managed by Terraform"
  private_key = tls_private_key.homelab.private_key_pem
}

# Create and validate server (automated)
resource "coolify_server" "homelab" {
  name             = "Homelab Server"
  description      = "Production server"
  ip               = "192.168.1.100"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.homelab.uuid
  instant_validate = true
  
  # Wait for validation - no manual UI clicks needed
  wait_for_validation = true
  
  timeouts {
    create = "5m"
  }
}

# Deploy database service (automated health check)
resource "coolify_service" "postgres" {
  name        = "PostgreSQL"
  server_uuid = coolify_server.homelab.uuid
  project_uuid = coolify_project.main.uuid
  environment_name = "production"
  
  instant_deploy = true
  
  # Wait for database to be healthy
  wait_for_deployment = true
  
  timeouts {
    create = "15m"
  }
  
  compose = <<EOF
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: secure_password
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
EOF
}

# Deploy API service (automated health check)
resource "coolify_service" "api" {
  name        = "API Service"
  server_uuid = coolify_server.homelab.uuid
  project_uuid = coolify_project.main.uuid
  environment_name = "production"
  
  instant_deploy = true
  
  # Wait for API to be healthy
  wait_for_deployment = true
  
  timeouts {
    create = "15m"
    update = "10m"  # Faster redeployments
  }
  
  compose = <<EOF
services:
  api:
    image: myorg/api:v1.0.0
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s
EOF
}

# Configure API with auto-restart
resource "coolify_service_envs" "api_config" {
  uuid = coolify_service.api.uuid
  
  # Automatically restart when configuration changes
  redeploy_on_change = true
  
  env {
    key   = "DATABASE_URL"
    value = "postgresql://postgres:secure_password@postgres:5432/mydb"
  }
  
  env {
    key   = "API_KEY"
    value = var.api_key
  }
  
  env {
    key   = "LOG_LEVEL"
    value = "info"
  }
}

# Output ready-to-use URLs
output "api_url" {
  value = "http://${coolify_server.homelab.ip}:8080"
  description = "API is deployed and healthy"
}
```

### Deployment Flow

1. **`terraform apply`** starts
2. Server created → automatically validates → continues when ready
3. PostgreSQL service deployed → waits for healthy → continues when ready  
4. API service deployed → waits for healthy → continues when ready
5. Environment variables set → service automatically restarts
6. **`terraform apply`** completes - everything is working

**Total manual steps: 0** ✨

## Troubleshooting

See the [Troubleshooting Guide](troubleshooting.md) for common issues and solutions.

## Best Practices Summary

1. ✅ **Always use `wait_for_validation` for production servers**
2. ✅ **Always use `wait_for_deployment` for critical services**
3. ✅ **Define healthchecks in Docker Compose** for accurate status detection
4. ✅ **Set appropriate timeouts** based on service complexity
5. ✅ **Use `redeploy_on_change` for runtime configuration**
6. ✅ **Chain dependent resources** to ensure proper ordering
7. ✅ **Test timeout values** with your actual infrastructure

## Performance Tips

The automation features add polling overhead but ensure correctness:

- **Server validation**: ~30 seconds to 3 minutes (depends on server)
- **Service deployment**: ~1-10 minutes (depends on image size and startup)
- **Redeploy**: ~30 seconds to 2 minutes (depends on service)

To optimize:
- Use `instant_validate = true` - starts validation immediately
- Use `instant_deploy = true` - starts deployment immediately  
- Set realistic but not excessive timeouts
- Cache Docker images on servers for faster deployments
