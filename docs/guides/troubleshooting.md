---
page_title: "Troubleshooting Guide"
description: |-
  Common issues and solutions when using the Coolify Terraform Provider
---

# Troubleshooting Guide

This guide covers common issues you might encounter when using the Coolify Terraform Provider and how to resolve them.

## Server Validation Issues

### Server Validation Timeout

**Symptom**: Server creation times out with "timeout waiting for server validation"

**Possible Causes**:
- Server is not reachable from Coolify instance (firewall, network configuration)
- SSH key authentication failing
- Server resources insufficient (disk space, memory)
- Port 22 is blocked or SSH is not running

**Solutions**:

1. **Increase timeout**:
```terraform
resource "coolify_server" "example" {
  # ... other config
  wait_for_validation = true
  
  timeouts {
    create = "10m"  # Increase from default 3m
  }
}
```

2. **Check server logs** in Coolify UI:
   - Navigate to Server → Validate
   - Review `validation_logs` for specific errors

3. **Verify network connectivity**:
   - Ensure server IP is reachable from Coolify instance
   - Check firewall rules allow SSH (port 22)
   - Verify SSH service is running: `systemctl status sshd`

4. **Test SSH manually**:
```bash
ssh -i /path/to/private_key user@server_ip
```

### Validation Fails Immediately

**Symptom**: Validation fails with error in logs

**Common Errors**:

- **"Permission denied (publickey)"**: SSH key mismatch or incorrect user
  - Verify `private_key_uuid` references correct key
  - Check `user` matches server's SSH user (usually `root`)
  
- **"Host key verification failed"**: Known_hosts conflict
  - Coolify should handle this automatically
  - If persists, check Coolify server's SSH configuration

## Service Deployment Issues

### Service Deployment Timeout

**Symptom**: Service creation times out with "timeout waiting for deployment"

**Possible Causes**:
- Docker image pull taking too long
- Service failing healthchecks
- Container startup errors
- Resource constraints (CPU, memory)

**Solutions**:

1. **Increase deployment timeout**:
```terraform
resource "coolify_service" "example" {
  # ... other config
  wait_for_deployment = true
  
  timeouts {
    create = "20m"  # Increase from default 10m
    update = "15m"
  }
}
```

2. **Check service status** in Coolify UI:
   - Navigate to Service → Logs
   - Review deployment logs and container logs

3. **Verify Docker Compose configuration**:
```terraform
resource "coolify_service" "example" {
  compose = <<EOF
services:
  app:
    image: "myapp:latest"
    healthcheck:  # Add healthcheck for better status detection
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
EOF
}
```

### Service Shows "degraded:unhealthy"

**Symptom**: Service deploys but status shows unhealthy

**Solutions**:

1. **Add or adjust healthcheck** in Docker Compose
2. **Increase healthcheck retries** to allow more time for startup
3. **Check application logs** for startup errors
4. **Verify environment variables** are correct

### Service Stuck in "deploying" State

**Symptom**: Service never reaches "running:healthy"

**Common Causes**:
- Missing or incorrect healthcheck configuration
- Application failing to start properly
- Port conflicts

**Debug Steps**:

1. Check actual container status: `docker ps -a`
2. Review container logs: `docker logs <container_name>`
3. Disable waiter temporarily to inspect service:
```terraform
resource "coolify_service" "example" {
  wait_for_deployment = false  # Deploy without waiting
  # ... rest of config
}
```

## Environment Variable Issues

### Service Not Restarting After Env Change

**Symptom**: Changed environment variables but service still using old values

**Solution**: Enable auto-redeploy:
```terraform
resource "coolify_service_envs" "example" {
  uuid = coolify_service.example.uuid
  
  redeploy_on_change = true  # Enable automatic restart
  
  env {
    key   = "MY_VAR"
    value = "new_value"
  }
}
```

### Redeploy Fails But Env Vars Updated

**Symptom**: Warning about redeploy failure but Terraform succeeds

**Explanation**: This is expected behavior - the provider updates environment variables successfully but issues a warning if the automatic redeployment request fails. The service will use new values on next manual restart.

**Solution**: If automatic redeployment is critical, check Coolify service status and redeploy manually if needed.

## Database Issues

### MariaDB/MongoDB/Redis/Clickhouse/KeyDB/Dragonfly Not Available

**Symptom**: "No resource named 'coolify_mariadb_database' available"

**Explanation**: These database types are currently **blocked due to incomplete Coolify OpenAPI specification**. The API endpoints return HTTP 200 instead of 201 and don't include response schemas.

**Workaround**: Use PostgreSQL or MySQL databases, which are fully supported:
```terraform
resource "coolify_postgresql_database" "db" {
  name              = "mydb"
  postgres_password = "secure_password"
  # ... other config
}

resource "coolify_mysql_database" "db" {
  name           = "mydb"
  mysql_password = "secure_password"
  # ... other config
}
```

**Future**: This will be resolved once Coolify's OpenAPI spec is updated. Track progress in the Coolify repository issues.

## General Tips

### Enable Debug Logging

Set environment variable for detailed logs:
```bash
export TF_LOG=DEBUG
terraform apply
```

### Check API Compatibility

Verify your Coolify version supports the features you're using:
- Server validation waiter: Any v4 version
- Service deployment waiter: Any v4 version  
- Redeploy on change: Any v4 version
- Database support: PostgreSQL and MySQL only

### Verify API Token

Test API access:
```bash
curl -H "Authorization: Bearer $COOLIFY_TOKEN" \
  https://your-coolify.com/api/v1/teams
```

### State Management

If resources get out of sync:
```bash
# Refresh state from actual infrastructure
terraform refresh

# Or import existing resources
terraform import coolify_server.example <server_uuid>
```

## Getting Help

If you encounter issues not covered here:

1. **Check the logs**: Enable `TF_LOG=DEBUG` for detailed output
2. **Review Coolify logs**: Check the Coolify UI for service/server specific errors
3. **Search existing issues**: Look for similar problems in GitHub issues
4. **Create an issue**: Provide:
   - Terraform version
   - Provider version
   - Coolify version
   - Full error message (redact sensitive data)
   - Minimal reproduction case

## Contributing

Found a solution not listed here? Please contribute by:
1. Opening a PR to update this guide
2. Sharing your experience in GitHub Discussions
3. Helping other users in issues
