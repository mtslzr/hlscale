module "dynamo" {
  source = "git::https://github.com/mtslzr/terraform.git//modules/dynamo?ref=0.1.2"

  hash_key       = "id"
  hash_key_type  = "S"
  project_name   = local.project_name
  range_key      = "start"
  range_key_type = "N"
  stream         = true
  table_name     = "exams"
}

resource "aws_lambda_event_source_mapping" "dynamo" {
  depends_on             = [module.lambda]
  event_source_arn       = module.dynamo.stream_arn
  function_name          = module.lambda.function_name

  maximum_retry_attempts = 0
  starting_position      = "LATEST"
}