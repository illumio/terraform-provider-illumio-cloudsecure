resource "illumio-cloudsecure_aws_flow_logs_s3_bucket_source" "flow_log_source" {
  account_id  = "812713887999"
  bucket_name = "my-vpc-flow-logs-bucket"
  path_prefix = "AWSLogs/812713887999/vpcflowlogs/"
}
