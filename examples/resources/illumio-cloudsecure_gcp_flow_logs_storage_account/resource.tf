/*
  ⚠️ Internal Use Only: This resource is intended strictly for use within internal modules that
  directly manage GCP storage account registration for flow log ingestion in Illumio CloudSecure.

  Do not use this resource directly unless you are building or extending internal provisioning logic.

  For production and standard use cases, leverage the officially supported Illumio CloudSecure modules,
  which abstract permissions and configuration for secure, consistent deployment.
*/

resource "illumio-cloudsecure_gcp_flow_logs_storage_account" "managed_flow_log_storage_account" {
  project_id                  = "my-gcp-project-id"
  storage_account_resource_id = "projects/my-gcp-project-id/buckets/my-flow-logs-bucket"
}
