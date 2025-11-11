# Example: Public Git Repository Application

resource "coolify_application" "public_example" {
  source_type = "public"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  git_repository = "https://github.com/coollabsio/coolify"
  git_branch     = "main"
  build_pack     = "nixpacks"
  ports_exposes  = "80"

  name        = "Public App Example"
  description = "Application from public Git repository"
}

