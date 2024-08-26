resource "illumio-cloudsecure_k8s_cluster_onboarding_credential" "example" {
  illumio_region = "aws-us-west-2"
  name           = "Dev clusters AWS US"

  # Optional attributes
  description = "Credential to onboard EKS dev clusters in AWS US regions"
}

output "example_client_id" {
  value       = illumio-cloudsecure_k8s_cluster_onboarding_credential.example.client_id
  description = "The client_id to use to onboard k8s clusters."
}

output "example_client_secret" {
  value       = illumio-cloudsecure_k8s_cluster_onboarding_credential.example.client_secret
  description = "The client_secret to use to onboard k8s clusters."
  sensitive   = true
}
