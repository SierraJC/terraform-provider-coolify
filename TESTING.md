# Test Coverage Summary

## Tests Added ✅

### 1. Wait Package Tests (`internal/wait/wait_test.go`)
**Coverage: 100%**

Tests added:
- `TestWaitForCondition_ImmediateSuccess` - Condition met on first check
- `TestWaitForCondition_SuccessAfterRetries` - Condition met after polling
- `TestWaitForCondition_ErrorFromCheck` - Error handling from check function
- `TestWaitForCondition_ContextTimeout` - Context deadline exceeded
- `TestWaitForCondition_ContextCanceled` - Context cancellation
- `TestWaitForCondition_ErrorOnSecondCall` - Transient error handling
- `TestTimeoutError` - Timeout error message formatting
- `TestWaitForCondition_Timing` - Verify polling interval behavior

### 2. Environment Resource Tests (`internal/service/environment_resource_test.go`)

Tests added:
- `TestAccEnvironmentResource` - Full CRUD lifecycle test
  - Create and Read operations
  - Import by UUID
  - Import by name
  - Update operations (replacement)
  - Delete (implicit)
- `TestAccEnvironmentResource_InvalidImportId` - Error handling for invalid import format
- `TestAccEnvironmentResource_MinimalConfig` - Minimal required fields

## Additional Tests to Consider

### Critical (Should Add Before PR)

#### 1. Server Validation Waiter Tests
**Why:** Core feature that modifies server resource behavior

Suggested tests:
```go
// internal/service/server_resource_test.go

func TestAccServerResource_WaitForValidation(t *testing.T) {
    // Test that wait_for_validation polls until server is ready
}

func TestAccServerResource_WaitForValidationTimeout(t *testing.T) {
    // Test timeout behavior with unreachable server
}
```

**Unit tests** (don't require real Coolify instance):
```go
// internal/service/server_resource_test.go

func TestServerResource_ValidationWaitLogic(t *testing.T) {
    // Test the conditional logic for when to wait
}
```

#### 2. Service Deployment Waiter Tests 
**Why:** Core feature that modifies service resource behavior

Suggested tests:
```go
// internal/service/service/resource_service_test.go

func TestAccServiceResource_WaitForDeployment(t *testing.T) {
    // Test wait_for_deployment polls until service is healthy
}

func TestAccServiceResource_WaitForDeploymentOnUpdate(t *testing.T) {
    // Test waiting on update operations
}
```

#### 3. Service Auto-Redeploy Tests
**Why:** Core feature that triggers deployments

Suggested tests:
```go
// internal/service/service_envs_resource_test.go

func TestAccServiceEnvsResource_RedeployOnChange(t *testing.T) {
    // Test that changing envs triggers redeploy when enabled
}

func TestAccServiceEnvsResource_NoRedeployWhenDisabled(t *testing.T) {
    // Test that envs can change without redeploy when disabled
}
```

### Nice to Have (Can Add Later)

#### 4. Integration Tests
Test combinations of features:
```go
func TestAccServerToServiceFlow(t *testing.T) {
    // Create server with wait_for_validation
    // Create service on validated server with wait_for_deployment
    // Verify end-to-end automation
}
```

#### 5. Edge Cases
```go
func TestAccServerResource_InstantValidateWithWait(t *testing.T) {
    // Test interaction between instant_validate and wait_for_validation
}

func TestAccServiceEnvsResource_RedeployWithWait(t *testing.T) {
    // Test redeploy_on_change + wait_for_deployment together
}
```

#### 6. Error Handling Tests
```go
func TestAccServerResource_ValidationFailure(t *testing.T) {
    // Test graceful handling when validation never succeeds
}
```

## Test Execution

### Run All Tests
```bash
make test
```

### Run Specific Package
```bash
go test ./internal/wait/ -v -cover
go test ./internal/service/ -v -run TestAccEnvironment
```

### Run with Coverage Report
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Acceptance Tests (Requires COOLIFY_TOKEN)
```bash
export COOLIFY_TOKEN="your-token"
make testacc
```

## Test Coverage Statistics

After adding these tests:

- `internal/wait`: **100%** ✅
- `internal/service/environment_resource`: **Acceptance tests added** ✅
- `internal/service/server_resource`: Needs waiter tests ⚠️
- `internal/service/service`: Needs waiter tests ⚠️
- `internal/service/service_envs_resource`: Needs redeploy tests ⚠️

## Recommendation for PR

**Before submitting PR to upstream:**

1. ✅ **Required** - Already done
   - Wait package unit tests (100%)
   - Environment resource acceptance tests

2. ⚠️ **Highly Recommended** - Should add
   - Server validation waiter unit tests (mocked)
   - Service deployment waiter unit tests (mocked)
   - Service redeploy logic unit tests

3. 📋 **Optional** - Can be added later or in follow-up PR
   - Full acceptance tests requiring live Coolify instance
   - Integration tests
   - Edge case coverage

## How to Add More Tests Safely

### For Features with External Dependencies (Server/Service)

**Option 1: Mock the API Client**
```go
// Create a mock that implements the API interface
type mockClient struct {
    getServerFn func(ctx context.Context, uuid string) (*api.Server, error)
}

func TestServerValidationLogic(t *testing.T) {
    mock := &mockClient{
        getServerFn: func(ctx context.Context, uuid string) (*api.Server, error) {
            // Return controlled responses
        },
    }
    // Test your logic
}
```

**Option 2: Table-Driven Tests**
```go
func TestServerValidationConditions(t *testing.T) {
    tests := []struct{
        name       string
        isReachable bool
        isUsable   bool
        shouldWait bool
    }{
        {"both true", true, true, false},
        {"reachable only", true, false, true},
        // ...
    }
    // Run tests
}
```

### For Acceptance Tests

Create test helpers in `internal/acctest/`:
```go
// acctest/helpers.go
func CreateTestServer(t *testing.T, config ServerConfig) string {
    // Reusable server creation logic
}
```

## CI/CD Considerations

The GitHub Actions workflow already runs:
- Unit tests on every PR
- Coverage reporting (when CODECOV_TOKEN set)
- Code generation verification

New tests will automatically run in CI, ensuring code quality.
