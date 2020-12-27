data "aws_iam_policy_document" "cwevents" {
  statement {
    effect = "Allow"
    actions = [
      "lambda:InvokeFunction"
    ]
    resources = [
      module.lambda.lambda_arn
    ]
  }
}

resource "aws_iam_policy" "cwevents" {
  name        = "${local.project_name}-cwevents"
  description = "Allow CloudWatch events to trigger Lambda."
  policy      = data.aws_iam_policy_document.lambda.json
}

resource "aws_iam_role" "cwevents" {
  name               = "${local.project_name}-cwevents"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "events.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cwevents" {
  policy_arn = aws_iam_policy.cwevents.arn
  role       = aws_iam_role.cwevents.name
}

data "aws_iam_policy_document" "lambda" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem"
    ]
    resources = [
      module.dynamo.table_arn
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:DescribeStream",
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:ListStreams"
    ]
    resources = [
      module.dynamo.stream_arn
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "events:PutRule",
      "events:PutTargets"
    ]
    resources = ["*"]
  }
  statement {
    effect = "Allow"
    actions = [
      "iam:PassRole"
    ]
    resources = [
      "arn:aws:iam::645714156459:role/${local.project_name}-cwevents"
    ]
  }
}

resource "aws_iam_policy" "lambda" {
  name        = local.project_name
  description = "Allow ${local.project_name} access to DynamoDB."
  policy      = data.aws_iam_policy_document.lambda.json
}