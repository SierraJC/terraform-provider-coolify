# Example: Dockerfile Application

resource "coolify_application" "dockerfile_example" {
  source_type = "dockerfile"

  project_uuid     = "your-project-uuid"
  server_uuid      = "your-server-uuid"
  environment_name = "production"

  dockerfile = <<EOF
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
EOF

  name        = "Dockerfile Example"
  description = "Application built from Dockerfile"
}

