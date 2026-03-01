## Summary

<!-- Please provide a clear and concise description of your changes -->

## Type of Change

<!-- Check all that apply -->

- [ ] 🐛 Bug fix (non-breaking change that fixes an issue)
- [ ] ✨ New feature (non-breaking change that adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Refactoring (no functional changes)
- [ ] 🧪 Test improvement
- [ ] 🏗️ Build/CI improvement
- [ ] 🧹 Chore (dependency updates, etc.)

## Related Issue

<!-- Link to the issue this PR addresses -->
Closes #(issue number)

<!-- Or if it doesn't close an issue -->
<!-- Related to #(issue number) -->

## Changes Made

<!-- Describe the changes you made in detail -->

-
-
-

## Testing

<!-- Describe how you tested your changes -->

### Unit Tests

- [ ] I have added/updated unit tests
- [ ] All unit tests pass (`make test`)

### Acceptance Tests (if applicable)

- [ ] I have added/updated acceptance tests
- [ ] I have tested against a live Coolify instance

### Manual Testing Steps

<!-- Describe manual testing you performed -->

1.
2.
3.

**Tested with:**
- Terraform version:
- Provider version:
- Coolify version:
- Go version:

## Documentation

<!-- Check all that apply -->

- [ ] I have updated the CHANGELOG.md
- [ ] I have updated/added documentation in `docs/`
- [ ] I have updated/added examples in `examples/`
- [ ] Documentation is auto-generated (resource/data source docs)
- [ ] No documentation changes needed

## Code Quality

<!-- Ensure these are checked before submitting -->

- [ ] My code follows the project's coding conventions
- [ ] I have run `go fmt` on my code
- [ ] I have run `go vet` and addressed any issues
- [ ] I have run `make test` and all tests pass
- [ ] I have added comments for complex logic
- [ ] My changes generate no new warnings

## Breaking Changes

<!-- If this is a breaking change, describe the impact and migration path -->

**Impact:**
<!-- Describe what will break -->

**Migration Path:**
<!-- How should users update their code? -->

```hcl
# Before
resource "coolify_example" "old" {
  # old way
}

# After
resource "coolify_example" "new" {
  # new way
}
```

## Screenshots (if applicable)

<!-- Add screenshots to help explain your changes -->

## Checklist

<!-- Final checks before submitting -->

- [ ] I have read the [CONTRIBUTING.md](CONTRIBUTING.md) guide
- [ ] My commits follow [Conventional Commits](https://www.conventionalcommits.org/) format
- [ ] I have updated the CHANGELOG.md (unreleased section)
- [ ] I have tested my changes thoroughly
- [ ] I have considered backward compatibility
- [ ] My code is ready for review

## Additional Notes

<!-- Any additional information reviewers should know -->
