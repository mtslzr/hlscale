terraform {
  backend "s3" {
    bucket         = "mtslzr-terraform"
    dynamodb_table = "terraform-lock"
    encrypt        = true
    key            = "hlstate/app/terraform.tfstate"
    region         = "us-east-1"
  }
}

locals {
  project_name = "hlscale"
}