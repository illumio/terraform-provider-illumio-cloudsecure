resource "illumio-cloudsecure_policy" "example" {
  name        = "Azure Network Policy"
  description = "Allow traffic between Azure network resources"
}

data "azurerm_subscription" "current" {}

data "azurerm_subnet" "destination" {
  name                 = "destination-subnet"
  resource_group_name  = "my-rg"
  virtual_network_name = "destination-vnet"
}

resource "illumio-cloudsecure_policy_version" "subscription_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from entire subscription to a specific subnet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            org_selector = {
              subscriptions = [
                { id = data.azurerm_subscription.current.subscription_id }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            network = {
              subnets = [
                { id = data.azurerm_subnet.destination.id }
              ]
            }
          }
        }
      }

      port_ranges = [
        {
          protocol  = "TCP"
          from_port = 443
          to_port   = 443
        }
      ]
    }
  ]
}

resource "illumio-cloudsecure_policy_provision" "example" {
  policy_version_id = illumio-cloudsecure_policy_version.subscription_to_subnet.id
}
