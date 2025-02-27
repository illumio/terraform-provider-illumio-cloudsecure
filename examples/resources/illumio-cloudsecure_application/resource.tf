data "azurerm_subscription" "current" {}
data "aws_caller_identity" "current" {}

# Create a deployment and an application

resource "illumio-cloudsecure_deployment" "test_deployment" {
  name                   = "Production"
  description            = "Production deployment"
  azure_subscription_ids = [data.azurerm_subscription.current.id]
  aws_account_ids        = [data.aws_caller_identity.current.account_id]
}


resource "illumio-cloudsecure_application" "test_application" {
  name          = "MyApplication"
  description   = "My example application"
  deployment_id = illumio-cloudsecure_deployment.test_deployment.id
}
