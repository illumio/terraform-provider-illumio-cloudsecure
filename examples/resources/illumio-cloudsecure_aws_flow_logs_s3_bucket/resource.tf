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
  ✅ Recommended Approach: Use the official Illumio CloudSecure module to register one or more AWS S3 buckets
  that store VPC flow logs with Illumio CloudSecure.

  This module manages the necessary permissions and resource configurations, aligning with
  best practices for secure and scalable integration.

  Module example:
  https://github.com/illumio/terraform-illumio-cloudsecure/tree/main/examples/aws_flow_logs_s3_buckets
*/