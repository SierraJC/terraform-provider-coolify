resource "coolify_mysql_database" "example" {
  name        = "Example Terraformed Database"
  description = "Managed by Terraform"

  server_uuid      = "rg8ks8c"
  project_uuid     = "uoswco88w8swo40k48o8kcwk"
  environment_name = "production"

  image               = "mysql:8"
  mysql_database      = "app"
  mysql_user          = "user"
  mysql_password      = "hunter12"
  mysql_root_password = "4-8-15-16-23-42"

  instant_deploy = false
}
