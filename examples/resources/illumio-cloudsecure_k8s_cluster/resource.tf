resource "illumio-cloudsecure_k8s_cluster" "example" {
  client_id       = "xxxxx"
  client_secret   = "xxxxx"
  illumio_region = "aws-us-west-2"

  # Optional attributes
  log_level = "debug"
}

output "example_client_id" {
  value       = illumio-cloudsecure_k8s_cluster.example.client_id
  description = "The service account client_id used."
}

output "example_client_secret" {
  value       = illumio-cloudsecure_k8s_cluster.example.client_secret
  description = "The service account client_secret used."
  sensitive   = true
}
