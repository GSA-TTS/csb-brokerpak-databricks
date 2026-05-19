resource "databricks_token" "binding_token" {
  comment          = "CSB model serving binding token for ${var.endpoint_name}"
  lifetime_seconds = var.token_lifetime_seconds > 0 ? var.token_lifetime_seconds : null
}