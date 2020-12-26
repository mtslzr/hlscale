terraform {
  backend "s3" {
    bucket         = "mtslzr-terraform"
    dynamodb_table = "terraform-lock"
    encrypt        = true
    key            = "hlstate/ci/terraform.tfstate"
    region         = "us-east-1"
  }
}

resource "aws_iam_role" "ci-role" {
  name               = "hlscale-ci"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::645714156459:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {}
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ci-role" {
  name   = "hlscale-ci"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:*",
        "cloudwatch:*",
        "dynamodb:*",
        "ec2:*",
        "iam:*",
        "lambda:*",
        "s3:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ci-role" {
  role       = aws_iam_role.ci-role.name
  policy_arn = aws_iam_policy.ci-role.arn
}