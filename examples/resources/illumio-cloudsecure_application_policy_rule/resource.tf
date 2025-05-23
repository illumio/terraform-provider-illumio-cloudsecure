resource "illumio-cloudsecure_application_policy_rule" "web_to_monitoring_allow_rule" {
  action         = "Allow"
  application_id = "d058a1ac-2e67-4f5b-8fb1-f0440eea6a67"
  description    = "Allow traffic from web servers to the monitoring role on specified TCP ports"

  from_ip_list_ids = ["49"]

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

  to_labels = [
    {
      key   = "role"
      value = "nagios"
    }
  ]

  to_ports = [
    {
      port_number = 90
      protocol    = "TCP"
    },
    {
      port_number = 443
      protocol    = "TCP"
    }
  ]
}
