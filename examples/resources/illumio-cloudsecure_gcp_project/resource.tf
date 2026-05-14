resource "illumio-cloudsecure_gcp_project" "managed_gcp_project" {
  project_id            = "my-dev-project"
  name                  = "Development GCP Project"
  mode                  = "ReadWrite"
  organization_id       = "organizations/123456789012"
  service_account_email = "cloudsecure@my-project.iam.gserviceaccount.com"
}
