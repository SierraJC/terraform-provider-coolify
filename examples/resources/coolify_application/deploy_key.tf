# Example: Private Deploy Key Application

resource "coolify_application" "deploy_key_example" {
  source_type = "private-deploy-key"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  private_key_uuid = "your-private-key-uuid"
  git_repository   = "git@github.com:your-org/your-private-repo.git"
  git_branch       = "main"
  build_pack       = "nixpacks"
  ports_exposes    = "80"

  name        = "Deploy Key Example"
  description = "Application from private repository via deploy key"
}

