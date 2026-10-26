resource "illumio-cloudsecure_policy" "example" {
  name        = "Cloud Network Policy"
  description = "Allow traffic between cloud network resources"
}

# --- Azure data sources ---

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

# Example 1: Azure subscription (org_selector) -> subnet
resource "illumio-cloudsecure_policy_version" "azure_subscription_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from entire Azure subscription to a specific subnet"

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

# Example 2: Azure VNet -> VNet
resource "illumio-cloudsecure_policy_version" "azure_vnet_to_vnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from source VNet to destination VNet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            network = {
              vnets = [
                {
                  subscription_id = data.azurerm_subscription.current.subscription_id
                  resource_group  = "my-rg"
                  id              = data.azurerm_virtual_network.source.name
                }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            network = {
              vnets = [
                {
                  subscription_id = data.azurerm_subscription.current.subscription_id
                  resource_group  = "my-rg"
                  id              = data.azurerm_virtual_network.destination.name
                }
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

# Example 3: Azure VNet -> subnet
resource "illumio-cloudsecure_policy_version" "azure_vnet_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from source VNet to destination subnet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            network = {
              vnets = [
                {
                  subscription_id = data.azurerm_subscription.current.subscription_id
                  resource_group  = "my-rg"
                  id              = data.azurerm_virtual_network.source.name
                }
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

# Example 4: Azure multiple subscriptions (org_selector) -> VNet
resource "illumio-cloudsecure_policy_version" "azure_multi_subscription_to_vnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from multiple Azure subscriptions to a specific VNet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          azure = {
            org_selector = {
              subscriptions = [
                { id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa" },
                { id = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb" },
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          azure = {
            network = {
              vnets = [
                {
                  subscription_id = data.azurerm_subscription.current.subscription_id
                  resource_group  = "my-rg"
                  id              = data.azurerm_virtual_network.destination.name
                }
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

# Example 5: AWS account (org_selector) -> subnet
resource "illumio-cloudsecure_policy_version" "aws_account_to_subnet" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from entire AWS account to a specific subnet"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          aws = {
            org_selector = {
              accounts = [
                { id = "123456789012" }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          aws = {
            network = {
              subnets = [
                {
                  id         = "subnet-0a1b2c3d4e5f67890"
                  account_id = "123456789012"
                  region     = "us-east-1"
                }
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

# Example 6: AWS VPC -> VPC
resource "illumio-cloudsecure_policy_version" "aws_vpc_to_vpc" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from source VPC to destination VPC"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          aws = {
            network = {
              vpcs = [
                {
                  id         = "vpc-0a1b2c3d4e5f67890"
                  account_id = "123456789012"
                  region     = "us-east-1"
                }
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          aws = {
            network = {
              vpcs = [
                {
                  id         = "vpc-0f9e8d7c6b5a43210"
                  account_id = "123456789012"
                  region     = "us-west-2"
                }
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

# Example 7: AWS multiple accounts (org_selector) -> VPC
resource "illumio-cloudsecure_policy_version" "aws_multi_account_to_vpc" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow HTTPS from multiple AWS accounts to a specific VPC"

  rules = [
    {
      action = "Allow"

      source = {
        cloud = {
          aws = {
            org_selector = {
              accounts = [
                { id = "111111111111" },
                { id = "222222222222" },
              ]
            }
          }
        }
      }

      destination = {
        cloud = {
          aws = {
            network = {
              vpcs = [
                {
                  id         = "vpc-0f9e8d7c6b5a43210"
                  account_id = "333333333333"
                  region     = "us-east-1"
                }
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
