variable "databricks_host" { type = string }
variable "databricks_token" {
  type      = string
  sensitive = true
}
variable "cluster_id" { type = string }
variable "token_lifetime_seconds" { type = number }
