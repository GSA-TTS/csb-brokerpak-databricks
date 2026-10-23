locals {
  cf_context      = jsondecode(var.cf_context_json)
  ttl_expires_at  = timeadd(timestamp(), format("%sh", var.ttl_hours))
  resource_tags   = merge(var.labels, {
    ManagedBy = "cloud-service-broker"
    TTLExpiry = local.ttl_expires_at
  })
  openai_base_url = replace(databricks_model_serving.endpoint.endpoint_url, "/invocations", "/v1")
}

output "endpoint_name" { value = databricks_model_serving.endpoint.name }
output "serving_endpoint_id" { value = databricks_model_serving.endpoint.serving_endpoint_id }
output "invocation_url" { value = databricks_model_serving.endpoint.endpoint_url }
output "openai_base_url" { value = local.openai_base_url }
output "databricks_host" { value = var.databricks_host }
output "served_entity_name" { value = var.served_entity_name }
output "entity_name" { value = var.entity_name }
output "entity_version" { value = var.entity_version }
output "workload_size" { value = var.workload_size }
output "workload_type" { value = var.workload_type }
output "scale_to_zero_enabled" { value = var.scale_to_zero_enabled }
output "budget_amount" { value = var.budget_amount }
output "ttl_expires_at" { value = local.ttl_expires_at }
output "resource_tags_json" { value = jsonencode(local.resource_tags) }
output "cf_provenance_json" {
  value = jsonencode({
    cf_organization_guid = try(local.cf_context.organization_guid, null)
    cf_organization_name = try(local.cf_context.organization_name, null)
    cf_space_guid        = try(local.cf_context.space_guid, null)
    cf_space_name        = try(local.cf_context.space_name, null)
  })
}