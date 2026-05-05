output "databricks_host" { value = var.databricks_host }
output "databricks_token" {
  value     = databricks_token.binding_token.token_value
  sensitive = true
}
output "cluster_id" { value = var.cluster_id }
