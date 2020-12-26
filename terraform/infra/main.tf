terraform {
  backend "s3" {
    bucket         = "mtslzr-terraform"
    dynamodb_table = "terraform-lock"
    encrypt        = true
    key            = "hlstate/infra/erraform.tfstate"
    region         = "us-east-1"
  }
}

resource "aws_launch_configuration" "launch" {
  name          = "hlscale-lc"
  image_id      = "ami-0885b1f6bd170450c"
  instance_type = "t2.micro"
}

resource "aws_autoscaling_group" "scale" {
  name             = "hlscale-asg"
  desired_capacity = 1
  min_size         = 1
  max_size         = 5

  launch_configuration = aws_launch_configuration.launch.name
  vpc_zone_identifier  = ["subnet-48b75874", "subnet-0356c666"]
}

resource "aws_s3_bucket" "website" {
  bucket = "hlscale-web"
  acl = "public-read"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::hlscale-web/*",
      "Principal": "*"
    }
  ]
}
EOF

  website {
    index_document = "index.html"
  }
}