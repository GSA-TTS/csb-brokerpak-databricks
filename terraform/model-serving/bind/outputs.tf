locals {
  cf_context = jsondecode(var.cf_context_json)
  binding_provenance = {
    cf_app_guid          = var.cf_app_guid
    cf_organization_guid = try(local.cf_context.organization_guid, null)
    cf_organization_name = try(local.cf_context.organization_name, null)
    cf_space_guid        = try(local.cf_context.space_guid, null)
    cf_space_name        = try(local.cf_context.space_name, null)
  }
  normalized_binding = {
    version            = "v1"
    provider           = "databricks"
    provisioner_family = "databricks_model_serving_token"
    connection_type    = "runtime"
    endpoint = {
      base_url    = var.openai_base_url
      region      = null
      api_version = "v1"
    }
    access = {
      mode       = "bearer_token"
      expires_at = null
    }
    grant = {
      kind                 = "service_principal"
      least_privilege_unit = "endpoint"
      allowed_models       = [var.entity_name]
    }
    credential = {
      format = "bearer_token"
      inline = {
        access_token = databricks_token.binding_token.token_value
      }
      secret_ref = null
    }
  }
}

output "databricks_host" { value = var.databricks_host }
output "databricks_token" {
  value     = databricks_token.binding_token.token_value
  sensitive = true
}
output "endpoint_name" { value = var.endpoint_name }
output "invocation_url" { value = var.invocation_url }
output "openai_base_url" { value = var.openai_base_url }
output "served_entity_name" { value = var.served_entity_name }
output "entity_name" { value = var.entity_name }
output "entity_version" { value = var.entity_version }
output "workload_size" { value = var.workload_size }
output "workload_type" { value = var.workload_type }
output "ttl_expires_at" { value = var.ttl_expires_at }
output "binding_provenance_json" { value = jsonencode(local.binding_provenance) }
output "resource_tags_json" { value = var.resource_tags_json }
output "normalized_binding_json" { value = jsonencode(local.normalized_binding) }