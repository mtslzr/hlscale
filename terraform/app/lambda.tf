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