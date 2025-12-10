/*
  ⚠️ Internal Use Only: This resource is intended strictly for use within internal modules that
  directly manage GCP Pub/Sub topic registration for flow log ingestion in Illumio CloudSecure.

  Do not use this resource directly unless you are building or extending internal provisioning logic.

  For production and standard use cases, leverage the officially supported Illumio CloudSecure modules,
  which abstract permissions and configuration for secure, consistent deployment.
*/

resource "illumio-cloudsecure_gcp_flow_logs_pubsub_topic" "managed_flow_log_pubsub_topic" {
  project_id      = "my-gcp-project-id"
  pubsub_topic_id = "projects/my-gcp-project-id/topics/my-flow-logs-topic"
}

/*
  ✅ Recommended Approach: Use the official Illumio CloudSecure module for GCP flow log onboarding.
  This module standardizes prerequisite IAM roles, project-level permissions,
  and resource registration for secure, repeatable adoption.
*/