terraform {
  backend "s3" {
    bucket         = "mtslzr-terraform"
    dynamodb_table = "terraform-lock"
    encrypt        = true
    key            = "hlstate/infra/erraform.tfstate"
    region         = "us-east-1"
  }
}

locals {
  home_ip      = "172.97.40.128"
  project_name = "hlscale"
}

resource "aws_launch_configuration" "launch" {
  name          = "${local.project_name}-lc"
  image_id      = "ami-0885b1f6bd170450c"
  instance_type = "t2.micro"
}

resource "aws_autoscaling_group" "scale" {
  name             = "${local.project_name}-asg"
  desired_capacity = 1
  min_size         = 1
  max_size         = 5

  launch_configuration = aws_launch_configuration.launch.name
  vpc_zone_identifier  = ["subnet-48b75874", "subnet-0356c666"]
}

resource "aws_s3_bucket" "website" {
  bucket = "${local.project_name}-web"
  acl    = "public-read"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:s3:::hlscale-web",
        "arn:aws:s3:::hlscale-web/*"
      ],
      "Principal": "*",
      "Condition": {
        "IpAddress": {
          "aws:SourceIp": "${local.home_ip}/32"
        }
      }
    }
  ]
}
EOF

  website {
    index_document = "index.html"
  }
}