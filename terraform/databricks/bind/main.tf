resource "databricks_token" "binding_token" {
  comment          = "CSB service binding token"
  lifetime_seconds = var.token_lifetime_seconds > 0 ? var.token_lifetime_seconds : null
}
