variable "databricks_host" { type = string }
variable "databricks_token" {
  type      = string
  sensitive = true
}
variable "endpoint_name" { type = string }
variable "invocation_url" { type = string }
variable "openai_base_url" { type = string }
variable "served_entity_name" { type = string }
variable "entity_name" { type = string }
variable "entity_version" { type = string }
variable "workload_size" { type = string }
variable "workload_type" { type = string }
variable "ttl_expires_at" { type = string }
variable "resource_tags_json" { type = string }
variable "cf_context_json" { type = string }
variable "cf_app_guid" { type = string }
variable "token_lifetime_seconds" { type = number }