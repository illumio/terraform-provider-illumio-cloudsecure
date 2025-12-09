/*
  ⚠️ Internal Use Only: This resource is intended strictly for use within internal modules that
  directly manage GCP Pub/Sub topic registration for flow log ingestion in Illumio CloudSecure.

  Do not use this resource directly unless you are building or extending internal provisioning logic.

  For production and standard use cases, leverage the officially supported Illumio CloudSecure modules,
  which abstract permissions and configuration for secure, consistent deployment.
*/

resource "illumio-cloudsecure_gcp_flow_logs_pub_sub" "managed_flow_log_pub_sub" {
  project_id                = "my-gcp-project-id"
  pub_sub_topic_resource_id = "projects/my-gcp-project-id/topics/my-flow-logs-topic"
}
