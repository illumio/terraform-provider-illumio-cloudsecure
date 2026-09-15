resource "illumio-cloudsecure_azure_flow_logs_storage_account_source" "flow_log_source" {
  subscription_id      = "00000000-0000-0000-0000-000000000000"
  storage_account_name = "myvpcflowlogs"
  container_name       = "insights-logs-networksecuritygroupflowevent"
  path_prefix          = "resourceId=/"
}
