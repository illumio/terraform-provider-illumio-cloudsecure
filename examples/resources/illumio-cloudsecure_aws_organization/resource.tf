resource "illumio-cloudsecure_aws_organization" "organization_example" {
  account_id       = "965208753613"
  organization_id  = "o-3eehyj6qk0"
  name             = "Test AWS Organization"
  role_arn         = "arn:aws:iam::965208753613:role/IllumioAccess"
  role_external_id = "eb287482f5824fab8a6988252d56eb6d"

  # Optional attributes
  mode = "ReadWrite"
}
