resource "illumio-cloudsecure_aws_account" "example" {
  account_id         = "123456789012"
  account_type       = "Organization"
  name               = "My AWS Account"
  service_account_id = "service-account-id"

  # Optional attributes
  disabled            = false
  excluded_regions    = ["us-west-1", "us-west-2"]
  excluded_subnet_ids = ["subnet-0123456789abcdef0", "subnet-abcdef0123456789"]
  excluded_vpc_ids    = ["vpc-0123456789abcdef0"]
  mode                = "ReadWrite"
}
