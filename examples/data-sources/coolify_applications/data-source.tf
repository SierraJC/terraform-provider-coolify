# Retrieve all applications
data "coolify_applications" "all" {}

# Retrieve applications with a specific UUID
data "coolify_applications" "filtered" {
  filter {
    name   = "uuid"
    values = ["abc123"]
  }
}

output "all" {
  value = data.coolify_applications.all.applications
}

output "filtered" {
  value = data.coolify_applications.filtered.applications
}
