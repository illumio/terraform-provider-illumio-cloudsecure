resource "illumio-cloudsecure_gcp_flow_logs_storage_bucket_source" "flow_log_source" {
  project_id  = "my-gcp-project"
  bucket_name = "my-vpc-flow-logs-bucket"
  path_prefix = "vpc-flow-logs/"
}
