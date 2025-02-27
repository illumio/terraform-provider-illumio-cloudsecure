data "azurerm_subscription" "current" {}
data "aws_caller_identity" "current" {}

resource "illumio-cloudsecure_deployment" "test_deployment" {
  name                   = "Production"
  description            = "Production deployment"
  azure_subscription_ids = [data.azurerm_subscription.current.id]
  aws_account_ids        = [data.aws_caller_identity.current.account_id]
}
