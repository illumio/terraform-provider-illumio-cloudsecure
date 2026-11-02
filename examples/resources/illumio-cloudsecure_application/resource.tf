data "azurerm_subscription" "current" {}
data "aws_caller_identity" "current" {}

# Define a deployment and an application

resource "illumio-cloudsecure_deployment" "test_deployment" {
  name                   = "Production"
  description            = "Production deployment"
  azure_subscription_ids = [data.azurerm_subscription.current.id]
  aws_account_ids        = [data.aws_caller_identity.current.account_id]
}

resource "illumio-cloudsecure_application" "test_application1" {
  name          = "MyApplication"
  description   = "My example application"
  deployment_id = illumio-cloudsecure_deployment.test_deployment.id
}

resource "illumio-cloudsecure_application" "test_application2" {
  name        = "AnotherApplication"
  description = "Another example application"
}
