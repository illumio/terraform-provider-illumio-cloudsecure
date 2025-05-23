---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "illumio-cloudsecure_k8s_cluster Resource - illumio-cloudsecure"
subcategory: ""
description: |-
  Manages the onboarding of a k8s cluster on CloudSecure in a specific Illumio Region.
---

# illumio-cloudsecure_k8s_cluster (Resource)

Manages the onboarding of a k8s cluster on CloudSecure in a specific Illumio Region.

## Example Usage

```terraform
resource "illumio-cloudsecure_k8s_cluster" "example" {
  illumio_region = "aws-us-west-2"

  # Optional attributes
  log_level = "Debug"
}

output "example_client_id" {
  value       = illumio-cloudsecure_k8s_cluster.example.client_id
  description = "The client_id to use to authenticate this k8s cluster."
}
output "example_client_secret" {
  value       = illumio-cloudsecure_k8s_cluster.example.client_secret
  description = "The client_secret to use to authenticate this k8s cluster."
  sensitive   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `illumio_region` (String) Illumio Region where the k8s cluster will be onboarded. An Illumio Region is a designated cloud region where the CloudSecure cloud-operator deployed in the k8s cluster connects after onboarding. Choose the Illumio Region nearest to the k8s cluster to maximize performance and security. Must be one of: `aws-ap-southeast-2`, `aws-eu-west-2`, `aws-us-west-2`, `aws-us-west-1`, `aws-eu-west-2`, `azure-us-east-2`, `azure-germany-west-central`, `azure-us-west-2`.

### Optional

- `log_level` (String) Verbosity of the logs produced by the CloudSecure k8s operator. Must be one of: `Debug`, `Info`, `Warn`, or `Error`.

### Read-Only

- `client_id` (String) Client identifier specific to the k8s cluster, used by CloudSecure's k8s operator to authenticate to CloudSecure, in combination with `client_secret`. Identical to `id`.
- `client_secret` (String, Sensitive) Client secret specific to the k8s cluster, used by CloudSecure's k8s operator to authenticate to CloudSecure, in combination with `client_id`.
- `id` (String) CloudSecure ID.
