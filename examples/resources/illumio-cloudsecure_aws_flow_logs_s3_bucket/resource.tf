/*
  ⚠️ Internal Use Only: Direct declaration of this resource is intended strictly for internal modules
  that manage low-level provisioning of AWS S3 flow log buckets within Illumio CloudSecure.

  Avoid using this resource directly unless you're developing or extending internal provisioning logic.

  For standard integrations, utilize the official Illumio CloudSecure module provided below,
  which abstracts the necessary configurations and promotes reusable, consistent infrastructure.

  Module documentation:
  https://registry.terraform.io/modules/illumio/cloudsecure/illumio/latest/submodules/aws_flow_logs_s3_buckets
*/

resource "illumio-cloudsecure_aws_flow_logs_s3_bucket" "flow_log_bucket" {
  account_id    = "812713887999"
  s3_bucket_arn = "arn:aws:s3:::flowlogbucket"
}

/*
  ✅ Recommended Approach: Use the following module to register one or more AWS S3 buckets
  that store VPC flow logs with Illumio CloudSecure.

  This module manages the necessary permissions and resource configurations, aligning with
  best practices for secure and scalable integration.
*/

module "aws_account_dev" {
  source                         = "illumio/cloudsecure/illumio//modules/aws_account"
  name                           = "Development AWS Account"
  illumio_cloudsecure_account_id = "158256226745"
  tags = {
    Name  = "CloudSecure Account"
    Owner = "Engineering"
  }
}

module "aws_flow_logs_s3_buckets" {
  source  = "illumio/cloudsecure/illumio//modules/aws_flow_logs_s3_buckets"
  version = "1.5.1"

  role_id = module.aws_account_dev.role_id

  s3_bucket_arns = [
    "arn:aws:s3:::flows-bucket-1",
    "arn:aws:s3:::flows-bucket-2",
    "arn:aws:s3:::flows-bucket-3",
    "arn:aws:s3:::flows-bucket-4/custom-path/first",
    "arn:aws:s3:::flows-bucket-5/custom-path/second",
  ]
}
