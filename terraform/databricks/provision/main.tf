resource "databricks_cluster" "cluster" {
  cluster_name            = var.cluster_name
  spark_version           = var.spark_version
  node_type_id            = var.node_type_id
  num_workers             = var.num_workers
  autotermination_minutes = var.autotermination_minutes

  custom_tags = var.labels

  lifecycle {
    prevent_destroy = true
  }
}
