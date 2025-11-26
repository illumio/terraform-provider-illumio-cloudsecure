resource "illumio-cloudsecure_application_policy_rule" "web_to_monitoring_allow_rule" {
  action         = "Allow"
  application_id = "d058a1ac-2e67-4f5b-8fb1-f0440eea6a67"
  description    = "Allow traffic from web servers to the monitoring role on specified TCP ports"

  from_ip_list_ids = ["c8427366-c5b8-49ed-a723-6fb3e9bb78a6"]

  from_labels = [
    {
      key   = "CostCenter"
      value = "1234"
    },
    {
      key   = "type"
      value = "server"
    }
  ]

  to_ip_list_ids = ["ea4436e1-436b-4b4d-8967-17cb946b53a5"]

  to_labels = [
    {
      key   = "role"
      value = "nagios"
    }
  ]

  to_port_ranges = [
    {
      from_port = 80
      to_port   = 80
      protocol  = "TCP"
    },
    {
      from_port = 443
      to_port   = 443
      protocol  = "TCP"
    },
  ]
}
