# Upstream Sync Workflow

How to sync changes from `SierraJC/terraform-provider-coolify` into our `wappsdev` fork.

## When to sync

- A meaningful upstream feature/fix you want
- Quarterly hygiene check (optional)
- NOT every renovate bot dep bump (we have our own)

## Procedure

```bash
cd terraform-provider-coolify
git fetch upstream
git checkout main
git merge upstream/main           # resolve conflicts manually if any

# Build to verify compatibility
go build -v ./...
go test ./...

# If clean: tag next minor version
git tag v1.X.0
git push origin main --tags
```

## Conflict resolution

`coolify_application_*` files were added by our fork (via PR #87). If upstream
later adds their own version (post-architectural-refactor), prefer keeping ours
and re-evaluating after compile + test pass.

## Provider registry republish

After tagging a new version, `.github/workflows/release.yml` builds via
GoReleaser and `.github/workflows/publish-registry.yml` regenerates the
`/docs/v1/` JSON manifests automatically.

## Rolling back a bad upstream sync

```bash
git revert -m 1 <merge-commit-sha>     # undoes the merge
git push origin main                    # consumers stay on previous version
```

Do NOT delete release tags — consumers may pin to specific versions.
