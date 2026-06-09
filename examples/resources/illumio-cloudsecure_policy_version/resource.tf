resource "illumio-cloudsecure_policy" "example" {
  name        = "Azure Network Policy"
  description = "Allow traffic between Azure network resources"
}

# Look up existing Azure resources
data "azurerm_subscription" "current" {}

data "azurerm_virtual_network" "source" {
  name                = "source-vnet"
  resource_group_name = "my-rg"
}

data "azurerm_virtual_network" "destination" {
  name                = "destination-vnet"
  resource_group_name = "my-rg"
}

data "azurerm_subnet" "destination" {
  name                 = "destination-subnet"
  resource_group_name  = "my-rg"
  virtual_network_name = "destination-vnet"
}

# Example 1: Whole subscription -> specific subnet
resource "illumio-cloudsecure_policy_version" "subscription_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from entire subscription to a specific subnet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
            # no network block = entire subscription
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
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

# Example 2: VNet -> VNet
resource "illumio-cloudsecure_policy_version" "vnet_to_vnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from source VNet to destination VNet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
            network = {
              vnets = [
                { id = data.azurerm_virtual_network.source.id }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
            network = {
              vnets = [
                { id = data.azurerm_virtual_network.destination.id }
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

# Example 3: VNet -> subnet
resource "illumio-cloudsecure_policy_version" "vnet_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from source VNet to destination subnet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
            network = {
              vnets = [
                { id = data.azurerm_virtual_network.source.id }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            subscription_id = data.azurerm_subscription.current.subscription_id
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
