resource "illumio-cloudsecure_aws_organization_account" "organization_account_example" {
  account_id                     = "987654321098"
  organization_master_account_id = "123456789012"
  role_arn                       = "arn:aws:iam::987654321098:role/IllumioAccess"
  role_external_id               = "eb287482f5824fab8a6988252d56eb6d"

  # Optional attributes
  disabled = false
}
