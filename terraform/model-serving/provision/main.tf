locals {
  endpoint_tags = [
    for key, value in var.labels : {
      key   = key
      value = tostring(value)
    }
  ]
}

resource "databricks_model_serving" "endpoint" {
  name        = var.endpoint_name
  description = var.description

  config {
    served_entities {
      name                  = var.served_entity_name
      entity_name           = var.entity_name
      entity_version        = var.entity_version
      workload_size         = var.workload_size
      workload_type         = var.workload_type
      scale_to_zero_enabled = var.scale_to_zero_enabled
    }

    traffic_config {
      routes {
        served_entity_name = var.served_entity_name
        traffic_percentage = 100
      }
    }
  }

  dynamic "tags" {
    for_each = local.endpoint_tags
    content {
      key   = tags.value.key
      value = tags.value.value
    }
  }
}