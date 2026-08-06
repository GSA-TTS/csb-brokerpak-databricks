variable "endpoint_name" { type = string }
variable "served_entity_name" { type = string }
variable "entity_name" { type = string }
variable "entity_version" { type = string }
variable "workload_size" { type = string }
variable "workload_type" { type = string }
variable "scale_to_zero_enabled" { type = bool }
variable "description" { type = string }
variable "budget_amount" { type = number }
variable "ttl_hours" { type = number }
variable "labels" { type = map(any) }
variable "cf_context_json" { type = string }
variable "databricks_host" { type = string }
variable "databricks_token" {
  type      = string
  sensitive = true
}