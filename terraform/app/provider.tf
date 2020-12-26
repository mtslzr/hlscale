provider "aws" {
  region = "us-east-1"
  version = "~> 3.0"

  assume_role {
    role_arn = "arn:aws:iam::645714156459:role/hlscale-ci"
    session_name = "terraform"
  }
}