resource "illumio-cloudsecure_aws_cloudtrail_s3_bucket" "cloudtrail_bucket" {
  account_id    = "812713887999"
  s3_bucket_arn = "arn:aws:s3:::cloudtrailbucket"
  s3_key_prefix = "org-trail"
}
