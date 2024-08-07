resource "illumio-cloudsecure_aws_account" "example" {
  account_id         = "812713887999"
  account_type       = "Organization"
  name               = "Test AWS Account"
  role_arn           = "arn:aws:iam::812713887999:role/IllumioAccess"
  service_account_id = "eb287482-f582-4fab-8a69-88252d56eb6d"

  # Optional attributes
  management_account_id = "965208753613"
  mode                  = "ReadWrite"
  organization_id       = "o-3eehyj6qk0"
}
