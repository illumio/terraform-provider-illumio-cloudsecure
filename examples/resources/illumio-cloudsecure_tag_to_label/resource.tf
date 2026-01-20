resource "illumio-cloudsecure_tag_to_label" "cloud_tag_environment" {
  name = "Environment"
  key  = "env"

  aws_tag_keys   = ["Environment", "Env"]
  azure_tag_keys = ["Environment"]
  gcp_tag_keys   = ["environment"]

  icon = {
    name             = "access"
    foreground_color = "#ffffff"
    background_color = "#1E90FF"
  }
}