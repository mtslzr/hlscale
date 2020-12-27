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
}

resource "aws_iam_policy" "lambda" {
  name        = local.project_name
  description = "Allow ${local.project_name} access to DynamoDB."
  policy      = data.aws_iam_policy_document.lambda.json
}

module "lambda" {
  depends_on = [aws_iam_policy.lambda]
  source     = "git::https://github.com/mtslzr/terraform.git//modules/lambda?ref=0.1.2"

  build_dir    = "${path.module}/../../build/lambda"
  description  = "Various utility functions for ${local.project_name}."
  filename     = "${path.module}/../../build/${local.project_name}.zip"
  handler      = local.project_name
  policy_arn   = aws_iam_policy.lambda.arn
  project_name = local.project_name
  runtime      = "go1.x"
}