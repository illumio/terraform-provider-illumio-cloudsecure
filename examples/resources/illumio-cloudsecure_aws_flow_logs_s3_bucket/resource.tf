resource "illumio-cloudsecure_aws_flow_logs_s3_bucket" "bucket_example" {
  account_id    = "812713887999"
  s3_bucket_arn = "arn:aws:s3:::exampleflowlogs"
}
