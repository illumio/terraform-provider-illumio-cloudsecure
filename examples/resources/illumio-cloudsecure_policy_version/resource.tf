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

# Example 8: Kubernetes default deny
#
# The "allow-dns" rule omits tier, so it is in the Admin tier.
# The "default-deny" rule denies all other traffic from the cluster's pods to any
# destination, inside or outside the cluster, in the Baseline tier.
# The "default-deny-ingress" rule denies traffic from any source to the
# cluster's pods, in the Baseline tier.
#
# policy_api on a clusters entry selects only clusters that provide that
# network policy API. Features only one API supports require it, so rules fail
# at plan rather than at provision: ingress from `any` or IP lists and service
# accounts need "cilium"; priority and the Baseline tier need "aws_vpc", except
# a Baseline default deny (a Deny to or from `any` on TCP 0-0 and UDP 0-0, all
# ports). Both APIs support the egress default deny, so "allow-dns" and
# "default-deny" omit policy_api and apply to every selected cluster. The
# ingress default deny needs "cilium", so it applies only to Cilium clusters.

variable "k8s_cluster_id" {
  type        = string
  description = "ID of the CloudSecure-onboarded Kubernetes cluster."
}

resource "illumio-cloudsecure_policy_version" "k8s_default_deny" {
  policy_id   = illumio-cloudsecure_policy.example.id
  description = "Allow DNS to kube-system, deny all other traffic"

  rules = {
    "allow-dns" = {
      action = "Allow"

      source = {
        k8s = {
          clusters           = [{ id = var.k8s_cluster_id }]
          namespace_selector = {}
          workload_selector  = {}
        }
      }

      destination = {
        k8s = {
          clusters = [{ id = var.k8s_cluster_id }]
          namespace_selector = {
            match_labels = [{ key = "kubernetes.io/metadata.name", value = "kube-system" }]
          }
          workload_selector = {
            match_labels = [{ key = "k8s-app", value = "kube-dns" }]
          }
        }
      }

      port_ranges = [
        { protocol = "UDP", from_port = 53, to_port = 53 },
        { protocol = "TCP", from_port = 53, to_port = 53 },
      ]
    }

    "default-deny" = {
      action = "Deny"

      source = {
        k8s = {
          clusters           = [{ id = var.k8s_cluster_id }]
          namespace_selector = {}
          workload_selector  = {}
          tier               = "Baseline"
        }
      }

      destination = {
        any = {}
      }

      port_ranges = [
        { protocol = "TCP", from_port = 0, to_port = 0 },
        { protocol = "UDP", from_port = 0, to_port = 0 },
      ]
    }

    "default-deny-ingress" = {
      action = "Deny"

      source = {
        any = {}
      }

      destination = {
        k8s = {
          clusters           = [{ id = var.k8s_cluster_id, policy_api = "cilium" }]
          namespace_selector = {}
          workload_selector  = {}
          tier               = "Baseline"
        }
      }

      port_ranges = [
        { protocol = "TCP", from_port = 0, to_port = 0 },
        { protocol = "UDP", from_port = 0, to_port = 0 },
      ]
    }
  }
}
