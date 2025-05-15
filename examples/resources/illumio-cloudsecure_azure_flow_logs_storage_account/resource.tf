/*
  ⚠️ Internal Use Only: This resource is intended strictly for use within internal modules that
  directly manage Azure storage account registration for flow log ingestion in Illumio CloudSecure.

  Do not use this resource directly unless you are building or extending internal provisioning logic.

  For production and standard use cases, leverage the officially supported Illumio CloudSecure modules,
  which abstract permissions and configuration for secure, consistent deployment.

  Module documentation:
  https://registry.terraform.io/modules/illumio/cloudsecure/illumio/latest/submodules/azure_flow_logs_storage_accounts
*/

resource "illumio-cloudsecure_azure_flow_logs_storage_account" "managed_flow_log_storage_account" {
  subscription_id             = "c219f111-9005-45d4-8bb3-4d50120d3ef2"
  storage_account_resource_id = "/subscriptions/randomids-d469-aghg-a4b4-asdsdasadas/resourceGroups/azrg-illumio/providers/Microsoft.Storage/storageAccounts/illumioazuretest"
}


/*
  ✅ Recommended: Use these modules to integrate Azure subscriptions and storage accounts with Illumio CloudSecure.

  These modules encapsulate authentication, role-based access, and flow log storage configuration,
  ensuring alignment with Illumio best practices and simplifying onboarding across environments.
*/

module "azure_subscription_dev" {
  source                       = "illumio/cloudsecure/illumio//modules/azure_subscription"
  version                      = "1.5.1"
  name                         = "Development Azure Subscription"
  mode                         = "ReadWrite"
  azure_secret_expiration_days = 365

  tags = [
    "Environment=Development",
    "Owner=John Doe",
  ]
}

module "azure_flow_log_storage_dev" {
  source                      = "illumio/cloudsecure/illumio//modules/azure_flow_logs_storage_accounts"
  version                     = "1.5.1"
  service_principal_client_id = module.azure_subscription_dev.service_principal_client_id

  storage_accounts = [
    {
      name                = "storage_bucket_name"
      resource_group_name = "storage_bucket_resource_group_name"
    },
  ]
}
