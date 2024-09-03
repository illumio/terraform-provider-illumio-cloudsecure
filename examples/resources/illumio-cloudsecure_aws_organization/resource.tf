resource "illumio-cloudsecure_aws_organization" "organization_example" {
  master_account_id = "965208753613"
  name              = "Test AWS Organization"

  # Optional attributes
  mode = "ReadWrite"
}
