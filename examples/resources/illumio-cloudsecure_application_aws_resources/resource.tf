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

resource "illumio-cloudsecure_application_aws_resources" "aws_arn_and_instance_resources" {
  application_id = illumio-cloudsecure_application.test_application.id
  account_id     = data.aws_caller_identity.current.account_id
  arns = [
    "arn:aws:ec2:us-east-1:600325505726:instance/i-0a1b2c3d4e5f67890",
    "arn:aws:rds:us-east-1:600325505726:cluster:my-database-cluster"
  ]
  aws_instances_ids = [
    "i-0a1b2c3d4e5f67890",
    "i-0b2c3d4e5f6a78901"
  ]
}


# Create an RDS cluster and add it to the application

resource "aws_rds_cluster" "example_cluster" {
  cluster_identifier = "my-tf-test-cluster"
  engine             = "aurora-mysql"
  master_username    = "admin"
  master_password    = "example-password"
}


resource "illumio-cloudsecure_application_aws_resources" "aws_rds_cluster_resources" {
  application_id = illumio-cloudsecure_application.test_application.id
  account_id     = data.aws_caller_identity.current.account_id
  aws_rds_cluster_ids = [
    aws_rds_cluster.example_cluster.id
  ]
}
