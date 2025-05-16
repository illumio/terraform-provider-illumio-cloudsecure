resource "illumio-cloudsecure_ip_list" "internal_subnet_range" {
  name        = "internal_subnet_range"
  description = "IP range for internal subnet access control"

  ip_ranges = [
    {
      exclusion       = false
      from_ip_address = "10.0.0.1"
      to_ip_address   = "10.0.0.255"
    }
  ]
}

resource "illumio-cloudsecure_ip_list" "corporate_networks" {
  name        = "corporate_networks"
  description = "List of allowed and excluded corporate IP blocks"

  ip_addresses = [
    {
      exclusion  = false
      ip_address = "10.0.0.0/16"
    },
    {
      exclusion  = true
      ip_address = "10.0.10.0/24"
    }
  ]
}
