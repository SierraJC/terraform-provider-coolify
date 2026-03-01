---
page_title: "Migration Guide - New Automation Features"
description: |-
  Upgrading guide for adopting operational automation features
---

# Migration Guide: Adopting Automation Features

This guide helps existing users adopt the new operational automation features without breaking existing deployments.

## Overview of New Features

This fork adds three major automation features:

1. **Server Validation Waiter** - `wait_for_validation` on `coolify_server`
2. **Service Deployment Waiter** - `wait_for_deployment` on `coolify_service`  
3. **Auto-Redeploy on Change** - `redeploy_on_change` on `coolify_service_envs`

All features are **opt-in** and backward compatible.

## Backward Compatibility

**Good News**: All new features are **optional** with safe defaults:

- `wait_for_validation = false` (default)
- `wait_for_deployment = false` (default)
- `redeploy_on_change = false` (default)

**Your existing configurations will continue to work unchanged.**

## Migration Strategies

### Strategy 1: Gradual Adoption (Recommended)

Add features incrementally to minimize risk:

#### Step 1: Add Server Validation (Low Risk)

```terraform
resource "coolify_server" "example" {
  # ... existing configuration unchanged
  
  # Add validation waiter
  wait_for_validation = true
  
  timeouts {
    create = "5m"
  }
}
```

**Impact**: Only affects new server creation or replacement.
**Test**: `terraform plan` - should show attribute additions, not resource replacement.

#### Step 2: Add Service Deployment Waiter (Medium Risk)

```terraform
resource "coolify_service" "example" {
  # ... existing configuration unchanged
  
  # Add deployment waiter
  wait_for_deployment = true
  
  timeouts {
    create = "15m"
    update = "10m"
  }
}
```

**Impact**: Only affects service redeployments, not existing running services.
**Test**: `terraform plan` - should show attribute additions only.

#### Step 3: Enable Auto-Redeploy (Highest Risk)

```terraform
resource "coolify_service_envs" "example" {
  # ... existing configuration unchanged
  
  # Add auto-redeploy
  redeploy_on_change = true
}
```

**Impact**: Next env var change will trigger automatic service restart.
**Test**: Make a non-critical env change first to verify behavior.

### Strategy 2: All At Once (Advanced Users)

If you're confident in your configuration:

```terraform
# Before
resource "coolify_server" "prod" {
  name             = "Production"
  ip               = "192.168.1.100"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.main.uuid
  instant_validate = false
}

resource "coolify_service" "app" {
  name        = "App"
  server_uuid = coolify_server.prod.uuid
  # ... rest of config
  instant_deploy = false
}

resource "coolify_service_envs" "app_config" {
  uuid = coolify_service.app.uuid
  # ... env vars
}

# After
resource "coolify_server" "prod" {
  name             = "Production"
  ip               = "192.168.1.100"
  port             = 22
  user             = "root"
  private_key_uuid = coolify_private_key.main.uuid
  instant_validate = true  # Changed
  
  # New
  wait_for_validation = true
  timeouts {
    create = "5m"
  }
}

resource "coolify_service" "app" {
  name        = "App"
  server_uuid = coolify_server.prod.uuid
  # ... rest of config
  instant_deploy = true  # Changed
  
  # New
  wait_for_deployment = true
  timeouts {
    create = "15m"
    update = "10m"
  }
}

resource "coolify_service_envs" "app_config" {
  uuid = coolify_service.app.uuid
  
  # New
  redeploy_on_change = true
  
  # ... env vars
}
```

## Testing Your Migration

### Pre-Migration Checklist

1. **Backup state**: 
```bash
cp terraform.tfstate terraform.tfstate.backup
```

2. **Review plan**:
```bash
terraform plan
```

Verify you see attribute additions, not resource replacements:
```
  ~ resource "coolify_server" "example" {
      + wait_for_validation = true
      + timeouts {
          + create = "5m"
        }
    }
```

3. **Test in non-production first**

### Post-Migration Validation

1. **Verify no unexpected changes**:
```bash
terraform plan  # Should show "No changes"
```

2. **Test automation on next deployment**:
   - Create a new server → should auto-validate
   - Deploy a service → should wait for health
   - Update env var → should auto-redeploy (if enabled)

## Common Migration Scenarios

### Scenario 1: Existing Infrastructure, New Services

You have existing servers and services, but want new deployments to use automation:

```terraform
# Existing server - don't change to avoid recreation
resource "coolify_server" "legacy" {
  name = "Legacy Server"
  # ... existing config, no waiters
}

# New server - use automation
resource "coolify_server" "new" {
  name = "New Server"
  # ... config
  wait_for_validation = true
  timeouts {
    create = "5m"
  }
}

# Existing service - don't change
resource "coolify_service" "legacy_app" {
  server_uuid = coolify_server.legacy.uuid
  # ... existing config, no waiters
}

# New service - use automation
resource "coolify_service" "new_app" {
  server_uuid = coolify_server.new.uuid
  # ... config
  wait_for_deployment = true
  timeouts {
    create = "15m"
  }
}
```

### Scenario 2: Convert Existing Resources Safely

To add waiters to existing resources without forcing recreation:

1. **Run plan to verify**:
```bash
terraform plan
```

2. **If you see resource replacement** (you shouldn't, but if you do):
```bash
# Instead of applying, fix the issue:
# - Check for other attribute changes
# - Verify no required attributes removed
# - Ensure values match current state
```

3. **Expected plan output**:
```hcl
# Should see updates, not replacements
~ resource "coolify_server" "example"
  + wait_for_validation = true
  # ...
```

### Scenario 3: Partial Adoption

Use automation only for critical services:

```terraform
# Production service - use all automation
resource "coolify_service" "prod_api" {
  name = "Production API"
  # ... config
  wait_for_deployment = true
  timeouts {
    create = "15m"
    update = "10m"
  }
}

resource "coolify_service_envs" "prod_api_config" {
  uuid = coolify_service.prod_api.uuid
  redeploy_on_change = true
  # ... env vars
}

# Development service - no automation for faster iteration
resource "coolify_service" "dev_api" {
  name = "Development API"
  # ... config
  # No waiters - manual validation in dev
}

resource "coolify_service_envs" "dev_api_config" {
  uuid = coolify_service.dev_api.uuid
  # redeploy_on_change defaults to false
  # ... env vars
}
```

## Rollback Procedure

If you need to revert:

1. **Remove automation features**:
```terraform
resource "coolify_server" "example" {
  # ... existing config
  
  # Comment out or remove
  # wait_for_validation = true
  # timeouts {
  #   create = "5m"
  # }
}
```

2. **Apply changes**:
```bash
terraform apply
```

3. **Restore from backup if needed**:
```bash
cp terraform.tfstate.backup terraform.tfstate
terraform refresh
```

## Performance Considerations

Adding waiters increases apply time but ensures correctness:

**Before** (no waiters):
- Server creation: ~5-10 seconds (API call only)
- Service deployment: ~10-30 seconds (API call only)
- Env var update: ~1-5 seconds (API call only)
- **Total**: ~16-45 seconds
- **Manual steps**: 3 (validate server, check service health, redeploy)

**After** (with waiters):
- Server creation + validation: ~1-3 minutes
- Service deployment + health check: ~2-10 minutes  
- Env var update + redeploy: ~1-2 minutes
- **Total**: ~4-15 minutes
- **Manual steps**: 0

**Trade-off**: Longer apply time vs. fully automated deployment.

## Best Practices for Migration

1. ✅ **Test in development first** - verify behavior before production
2. ✅ **Migrate gradually** - one feature at a time
3. ✅ **Keep backups** - save state before major changes
4. ✅ **Review plans** - always check `terraform plan` output
5. ✅ **Monitor first deploy** - watch the first automated deployment
6. ✅ **Document changes** - note which resources use automation
7. ✅ **Set appropriate timeouts** - based on your infrastructure

## Troubleshooting Migration Issues

### Issue: Plan Shows Resource Replacement

**Symptom**: `terraform plan` shows server/service will be destroyed and recreated

**Cause**: Usually means you changed a required attribute or schema changed

**Solution**:
```bash
# Check what changed
terraform plan -out=plan.out
terraform show plan.out

# If only automation attributes changed, this shouldn't happen
# Check for other changes you made
```

### Issue: First Apply Times Out

**Symptom**: First apply with waiters times out

**Solution**: Increase timeouts:
```terraform
timeouts {
  create = "10m"  # Increase from default
}
```

Review [Troubleshooting Guide](troubleshooting.md) for specific timeout issues.

### Issue: Auto-Redeploy Fails

**Symptom**: Warning about redeploy failure when updating env vars

**Explanation**: This is expected behavior - env vars are still updated successfully

**Solution**: 
- Check service status in Coolify UI
- Manually redeploy if needed
- Review service logs to diagnose redeploy failure

## Getting Help

If you encounter migration issues:

1. Check `terraform plan` output carefully
2. Review [Troubleshooting Guide](troubleshooting.md)
3. Search GitHub issues
4. Open a new issue with:
   - Before/after configurations
   - `terraform plan` output
   - Provider version
   - Coolify version

## Summary

- ✅ **All features are opt-in** - existing configs work unchanged
- ✅ **No forced resource replacements** - safe to add to existing resources
- ✅ **Gradual adoption recommended** - add features incrementally
- ✅ **Easy rollback** - remove attributes if needed
- ✅ **Test in dev first** - verify behavior before production
