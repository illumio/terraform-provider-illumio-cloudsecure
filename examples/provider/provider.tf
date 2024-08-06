# Configure the Illumio CloudSecure provider using the required_providers stanza.
terraform {
  required_providers {
    illumio-cloudsecure = {
      source  = "illumio/illumio-cloudsecure"
      version = "~> 1.0.2"
    }
  }
}

provider "illumio-cloudsecure" {
  client_id     = "my-access-id"
  client_secret = "my-secret-id"
}
