---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "illumio-cloudsecure_aws_organization Resource - illumio-cloudsecure"
subcategory: ""
description: |-
  Manages an AWS organization in CloudSecure.
---

# illumio-cloudsecure_aws_organization (Resource)

Manages an AWS organization in CloudSecure.

## Example Usage

```terraform
resource "illumio-cloudsecure_aws_organization" "organization_example" {
  master_account_id = "965208753613"
  organization_id   = "o-3eehyj6qk0"
  name              = "Test AWS Organization"
  role_arn          = "arn:aws:iam::965208753613:role/IllumioAccess"
  role_external_id  = "eb287482f5824fab8a6988252d56eb6d"

  # Optional attributes
  mode = "ReadWrite"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `master_account_id` (String) ID of the master account of the AWS organization.
- `name` (String) Display name.
- `organization_id` (String) AWS organization ID.
- `role_arn` (String) ARN of the AWS role to be assumed by CloudSecure to manage this account.
- `role_external_id` (String) External ID defined in the AWS role to authenticate CloudSecure when assuming that role.

### Optional

- `mode` (String) Access mode, must be `"ReadWrite"` (default) or `"Read"`.

### Read-Only

- `id` (String) CloudSecure ID.