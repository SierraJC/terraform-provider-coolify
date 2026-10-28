# Example: Private GitHub App Application

resource "coolify_application" "github_app_example" {
  source_type = "private-github-app"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  github_app_uuid  = "your-github-app-uuid"
  git_repository   = "https://github.com/your-org/your-private-repo"
  git_branch       = "main"
  build_pack       = "nixpacks"
  ports_exposes    = "80"

  name        = "GitHub App Example"
  description = "Application from private GitHub repository via GitHub App"
}

