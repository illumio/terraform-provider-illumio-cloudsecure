data "azurerm_subscription" "current" {}

# Define a deployment and an application

resource "illumio-cloudsecure_deployment" "test_deployment" {
  name                   = "Production"
  description            = "Production deployment"
  azure_subscription_ids = [data.azurerm_subscription.current.id]
}

resource "illumio-cloudsecure_application" "test_application" {
  name          = "MyApplication"
  description   = "My example application"
  deployment_id = illumio-cloudsecure_deployment.test_deployment.id
}


# Add existing Azure resources to the application
resource "illumio-cloudsecure_application_azure_resources" "azure_existing_resources" {
  application_id  = illumio-cloudsecure_application.test_application.id
  subscription_id = data.azurerm_subscription.current.id
  resource_ids = [
    "/subscriptions/01b93b91-c2f8-4702-82d0-69d2d4abfab5/resourceGroups/autorg1/providers/Microsoft.Network/networkInterfaces/auto-ohio-vm-1VMNic",
    "/subscriptions/01b93b91-c2f8-4702-82d0-69d2d4abfab5/resourceGroups/autorg21/providers/Microsoft.Compute/virtualMachines/auto-california-vm-18",
    "/subscriptions/01b93b91-c2f8-4702-82d0-69d2d4abfab5/resourceGroups/autorg22/providers/Microsoft.Network/network"
  ]
}

# Add new Azure resources to the application
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "illumio-cloudsecure_application_azure_resources" "azure_new_resources" {
  application_id  = illumio-cloudsecure_application.test_application.id
  subscription_id = data.azurerm_subscription.current.id
  resource_ids = [
    azurerm_resource_group.example.id
  ]
}