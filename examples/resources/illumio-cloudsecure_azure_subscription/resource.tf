/*
  ⚠️ Internal Use Only: This resource is intended solely for internal module development,
  where direct provisioning of Azure subscriptions is required for Illumio CloudSecure integration.

  Avoid using this resource directly in production or general use cases.

  Instead, use the officially supported `azure_subscription` module provided by Illumio,
  which abstracts authentication, permission scopes, and secure secret management.

  Module documentation:
  https://registry.terraform.io/modules/illumio/cloudsecure/illumio/latest/submodules/azure_subscription
*/

resource "illumio-cloudsecure_azure_subscription" "managed_azure_subscription" {
  client_id       = "ZDIASAD7RGBTESTJUPUJ"
  client_secret   = "iam12TsTe1s17h7M27e8REGw7oqGocKR2ZDveZsM"
  name            = "Development Azure Subscription"
  subscription_id = "6a879a4d-efdc-4b07-ad91-1919203356f5"
  tenant_id       = "de6b88d1-8289-4d5c-9453-f5c003e8dd51"

  # Optional configuration
  mode = "ReadWrite"
}


/*
  ✅ Recommended Approach: Use the official Illumio CloudSecure module to onboard Azure subscriptions with Illumio CloudSecure.

  This module simplifies authentication, secret handling, and role assignment,
  promoting maintainable and secure infrastructure-as-code.

  Module example:
  https://github.com/illumio/terraform-illumio-cloudsecure/tree/main/examples/azure_subscription
*/