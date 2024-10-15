resource "illumio-cloudsecure_aws_flow_logs_s3_bucket" "flow_log_bucket" {
  account_id    = "812713887999"
  s3_bucket_arn = "arn:aws:s3:::flowlogbucket"
}
