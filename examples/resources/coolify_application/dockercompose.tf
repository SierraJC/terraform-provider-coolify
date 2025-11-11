# Example: Docker Compose Application

resource "coolify_application" "dockercompose_example" {
  source_type = "dockercompose"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  docker_compose_raw = <<EOF
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: example
EOF

  name        = "Docker Compose Example"
  description = "Application from Docker Compose"
}

