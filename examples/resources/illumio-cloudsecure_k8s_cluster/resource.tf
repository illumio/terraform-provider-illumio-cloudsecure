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
