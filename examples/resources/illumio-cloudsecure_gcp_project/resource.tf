/*
  ⚠️ Internal Use Only: This resource is intended for internal module development
  where direct control over GCP project onboarding with Illumio CloudSecure is required.

  Avoid using this resource directly in production or typical deployments.

  Instead, leverage the official Illumio CloudSecure modules that encapsulate
  project onboarding workflows and recommended IAM configurations.
*/

resource "illumio-cloudsecure_gcp_project" "managed_gcp_project" {
  account_id            = "123456789012"
  name                  = "Development GCP Project"
  mode                  = "ReadWrite"
  organization_id       = "organizations/123456789012"
  service_account_email = "cloudsecure@my-project.iam.gserviceaccount.com"
}

/*
  ✅ Recommended Approach: Use the official Illumio CloudSecure module for GCP onboarding.
  This module standardizes prerequisite IAM roles, project-level permissions,
  and resource registration for secure, repeatable adoption.
*/
