resource "illumio-cloudsecure_gcp_flow_logs_pubsub_topic" "managed_flow_log_pubsub_topic" {
  project_id      = "my-gcp-project-id"
  pubsub_topic_id = "projects/my-gcp-project-id/topics/my-flow-logs-topic"
}
