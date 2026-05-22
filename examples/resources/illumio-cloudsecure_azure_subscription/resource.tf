resource "illumio-cloudsecure_azure_subscription" "managed_azure_subscription" {
  client_id       = "ZDIASAD7RGBTESTJUPUJ"
  client_secret   = "iam12TsTe1s17h7M27e8REGw7oqGocKR2ZDveZsM"
  name            = "Development Azure Subscription"
  subscription_id = "6a879a4d-efdc-4b07-ad91-1919203356f5"
  tenant_id       = "de6b88d1-8289-4d5c-9453-f5c003e8dd51"

  # Optional configuration
  mode = "ReadWrite"
}
