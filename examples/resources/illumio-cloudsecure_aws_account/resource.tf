resource "illumio-cloudsecure_aws_account" "managed_aws_account" {
  account_id       = "812713887999"
  name             = "Development AWS Account"
  role_arn         = "arn:aws:iam::812713887999:role/IllumioAccess"
  role_external_id = "eb287482f5824fab8a6988252d56eb6d"

  # Optional attributes
  mode            = "ReadWrite"
  organization_id = "o-3eehyj6qk0"
}
