/*
  ⚠️ Internal Use Only: This resource is intended for internal module development
  where fine-grained control over AWS account registration with Illumio CloudSecure is required.

  Do not use this resource directly in production or standard usage scenarios.

  Instead, leverage the official Illumio CloudSecure AWS Account module, which encapsulates
  permission handling, validation, and resource registration in a reusable interface.

  Module documentation:
  https://registry.terraform.io/modules/illumio/cloudsecure/illumio/latest/submodules/aws_account
*/

resource "illumio-cloudsecure_aws_account" "managed_aws_account" {
  account_id       = "812713887999"
  name             = "Development AWS Account"
  role_arn         = "arn:aws:iam::812713887999:role/IllumioAccess"
  role_external_id = "eb287482f5824fab8a6988252d56eb6d"

  # Optional attributes
  mode            = "ReadWrite"
  organization_id = "o-3eehyj6qk0"
}


/*
  ✅ Recommended: Use the official Illumio CloudSecure module to register an AWS account with Illumio CloudSecure.

  This module abstracts direct resource management, ensuring compliance with
  Illumio’s onboarding workflow and enforcing best practices for secure,
  scalable infrastructure integration.

  Module example:
  https://github.com/illumio/terraform-illumio-cloudsecure/tree/main/examples/aws_account
*/