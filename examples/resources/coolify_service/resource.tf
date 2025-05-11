resource "coolify_service" "example" {
  name        = "Example Terraformed Service"
  description = "Managed by Terraform"

  server_uuid      = "rg8ks8c"
  project_uuid     = "uoswco88w8swo40k48o8kcwk"
  environment_name = "production"
  destination_uuid = "kgso0w8"

  instant_deploy = false

  compose = <<EOF
services:
  whoami:
    image: "containous/whoami"
    container_name: "simple-service"
EOF

}
