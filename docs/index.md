---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "illumio-cloudsecure Provider"
subcategory: ""
description: |-
  A Provider for managing Illumio CloudSecure.
---

# illumio-cloudsecure Provider

A Provider for managing Illumio CloudSecure.

## Example Usage

```terraform
provider "scaffolding" {
  # example configuration here
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_token` (String, Sensitive) OAuth 2 access token used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.
- `api_endpoint` (String) CloudSecure Config API endpoint, defaults to dns:///cloud.illum.io:443.
- `client_id` (String) OAuth 2 client identifier used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.
- `client_secret` (String, Sensitive) OAuth 2 client secret used to authenticate against the CloudSecure Config API. Either client_id+client_secret or access_token must be specified.
- `disable_tls` (Boolean) Disables TLS for all all requests to the CloudSecure Token and Config API endpoints. TLS is enabled by default. Should only be used for testing the provider.
- `request_timeout` (String) Maximum duration of each API request, defaults to 10s.
- `token_endpoint` (String) CloudSecure OAuth 2 Token endpoint, defaults to https://cloud.illum.io/token.