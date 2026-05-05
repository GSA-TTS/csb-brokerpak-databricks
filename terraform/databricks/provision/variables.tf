variable "cluster_name" { type = string }
variable "spark_version" { type = string }
variable "node_type_id" { type = string }
variable "num_workers" { type = number }
variable "autotermination_minutes" { type = number }
variable "labels" { type = map(any) }
variable "databricks_host" { type = string }
variable "databricks_token" {
  type      = string
  sensitive = true
}
