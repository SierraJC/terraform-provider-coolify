# Example: Docker Image Application

resource "coolify_application" "dockerimage_example" {
  source_type = "dockerimage"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  docker_registry_image_name = "nginx"
  docker_registry_image_tag   = "alpine"
  ports_exposes              = "80"

  name        = "Docker Image Example"
  description = "Application from Docker image"
}

