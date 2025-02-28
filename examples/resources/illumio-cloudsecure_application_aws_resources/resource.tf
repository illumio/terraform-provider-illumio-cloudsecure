data "aws_caller_identity" "current" {}

# Define a deployment and an application

resource "illumio-cloudsecure_deployment" "test_deployment" {
  name            = "Production"
  description     = "Production deployment"
  aws_account_ids = [data.aws_caller_identity.current.account_id]
}

resource "illumio-cloudsecure_application" "test_application" {
  name          = "MyApplication"
  description   = "My example application"
  deployment_id = illumio-cloudsecure_deployment.test_deployment.id
}


# Add existing AWS resources to the application

resource "illumio-cloudsecure_application_aws_resources" "aws_security_group_resources" {
  application_id = illumio-cloudsecure_application.test_application.id
  account_id     = data.aws_caller_identity.current.account_id
  aws_security_group_ids = [
    "sg-021b2bc8d1f6b2dec",
    "sg-0742cd5a71ccbfc67"
  ]
}

resource "illumio-cloudsecure_application_aws_resources" "aws_arn_resources_and_internet_gateways" {
  application_id = illumio-cloudsecure_application.test_application.id
  account_id     = data.aws_caller_identity.current.account_id
  arns = [
    "arn:aws:kms:us-east-1:600325505726:key/6c44f35b-d3b5-4cef-9944-36b7df5d86c0",
    "arn:aws:s3bucketpolicy:::mys3bucketpolicy"
  ]
  aws_internet_gateway_ids = [
    "igw-0510c5ae3d648b857",
    "igw-0aa0dfc78b498845b"
  ]
}


# Create an S3 bucket and add it to the application

resource "aws_s3_bucket" "example_bucket" {
  bucket = "my-tf-test-bucket"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}


resource "illumio-cloudsecure_application_aws_resources" "aws_s3_bucket_resources" {
  application_id = illumio-cloudsecure_application.test_application.id
  account_id     = data.aws_caller_identity.current.account_id
  arns = [
    aws_s3_bucket.example_bucket.arn
  ]
}
