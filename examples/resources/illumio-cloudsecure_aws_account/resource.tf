resource "illumio-cloudsecure_aws_account" "account_example" {
  account_id                     = "812713887999"
  name                           = "Test AWS Account"
  role_arn                       = "arn:aws:iam::812713887999:role/IllumioAccess"
  role_external_id               = "eb287482f5824fab8a6988252d56eb6d"
  organization_master_account_id = "965208753613"
  organization_id                = "o-1234567890"
  access_mode                    = "ReadWrite"
}
