resource "coolify_service_envs" "example" {
  # Service UUID
  uuid = "i0800ok00gcww840kk8sok0s"

  env {
    key   = "key1"
    value = "value1"
  }

  env {
    key   = "key2"
    value = "value2"
  }
}
