resource "illumio-cloudsecure_aws_account" "example" {
  account_id         = "123456789012"
  account_type       = "Organization"
  name               = "My AWS Account"
  service_account_id = "service-account-id"

  # Optional attributes
  mode = "ReadWrite"
}
