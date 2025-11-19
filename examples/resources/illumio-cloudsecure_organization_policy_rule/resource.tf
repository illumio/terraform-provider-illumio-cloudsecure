resource "illumio-cloudsecure_organization_policy" "example" {
  name        = "org-policy-example"
  description = "Example organization policy"
  enabled     = true
}

resource "illumio-cloudsecure_ip_list" "example" {
  name        = "example"
  description = "IP range for example ip list"

  ip_ranges = [
    {
      exclusion       = false
      from_ip_address = "10.0.0.1"
      to_ip_address   = "10.0.0.255"
    }
  ]
}

resource "illumio-cloudsecure_organization_policy_rule" "example" {
  organization_policy_id = illumio-cloudsecure_organization_policy.example.id
  action                 = "Allow"
  description            = "Allow traffic from IP list to API role on ports 90-9000"

  from_ip_list_ids = [
    illumio-cloudsecure_ip_list.example.id
  ]

  to_labels = [
    {
      key   = "role"
      value = "API"
    }
  ]

  to_port_ranges = [
    {
      from_port = 90
      to_port   = 9000
      protocol  = "TCP"
    }
  ]
}