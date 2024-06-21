# Configure the Illumio CloudSecure provider using the required_providers stanza.
terraform {
  required_providers {
    illumio-cloudsecure = {
      source  = "illumio/illumio-cloudsecure"
      version = "~> 0.1"
    }
  }
}

# Onboard an AWS account
resource "illumio-cloudsecure_aws_account" "example" {
  # ...
}
