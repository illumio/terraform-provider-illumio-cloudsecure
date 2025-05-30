# Configure the Illumio CloudSecure provider using the required_providers stanza.
terraform {
  required_providers {
    illumio-cloudsecure = {
      source  = "illumio/illumio-cloudsecure"
      version = "~> 1.5.3"
    }
  }
}

provider "illumio-cloudsecure" {
  client_id     = "my-client-id"
  client_secret = "my-secret-id"
}
