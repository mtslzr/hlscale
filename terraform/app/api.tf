data "template_file" "api-swagger" {
  depends_on = [module.lambda]
  template   = file("api.yaml")
}

resource "aws_api_gateway_rest_api" "api" {
  depends_on = [module.lambda]
  name       = "${local.project_name}-api"
  body       = data.template_file.api-swagger.rendered
}

resource "aws_api_gateway_deployment" "api" {
  rest_api_id       = aws_api_gateway_rest_api.api.id
  stage_name        = "api"
  stage_description = "Deployed at ${timestamp()}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lambda_permission" "api" {
  action        = "lambda:InvokeFunction"
  function_name = local.project_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}