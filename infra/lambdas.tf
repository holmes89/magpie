module "lambda_label" {
  source     = "cloudposse/label/null"
  version    = "0.25.0"
  context    = module.this.context
  attributes = ["lambda"]
  enabled    = module.this.enabled
}

data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "assume" {
  count = module.this.enabled ? 1 : 0
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  count              = module.this.enabled ? 1 : 0
  name               = module.lambda_label.id
  assume_role_policy = join("", data.aws_iam_policy_document.assume.*.json)
  tags               = module.lambda_label.tags
}

data "aws_iam_policy_document" "lambda" { #Should I break it up?
  count = module.this.enabled ? 1 : 0
  # Dynamo Connection
  statement {
    actions = [
      "dynamodb:List*",
      "dynamodb:DescribeReservedCapacity*",
      "dynamodb:DescribeLimits",
      "dynamodb:DescribeTimeToLive"
    ]
    resources = ["*"]
  }
  statement {
    actions = [
      "dynamodb:BatchGet*",
      "dynamodb:DescribeStream",
      "dynamodb:DescribeTable",
      "dynamodb:Get*",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWrite*",
      "dynamodb:CreateTable",
      "dynamodb:Delete*",
      "dynamodb:Update*",
      "dynamodb:PutItem"
    ]
    resources = ["arn:aws:dynamodb:${data.aws_region.current.name}:*:table/${aws_dynamodb_table.magpie_table.name}"]
  }

}

resource "aws_iam_policy" "lambda" {
  count  = module.this.enabled ? 1 : 0
  name   = module.lambda_label.id
  policy = join("", data.aws_iam_policy_document.lambda.*.json)
}

resource "aws_iam_role_policy_attachment" "lambda" {
  count      = module.this.enabled ? 1 : 0
  role       = join("", aws_iam_role.lambda.*.name)
  policy_arn = join("", aws_iam_policy.lambda.*.arn)
}

resource "aws_iam_role_policy_attachment" "lambda" {
  count      = module.this.enabled ? 1 : 0
  role       = join("", aws_iam_role.lambda.*.name)
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

module "api_label" {
  source     = "cloudposse/label/null"
  version    = "0.25.0"
  context    = module.lambda_label.context
  attributes = ["api"]
  enabled    = module.this.enabled
}

resource "aws_lambda_function" "api" {
  count         = module.this.enabled ? 1 : 0
  function_name = module.api_label.id
  tags          = module.api_label.tags
  filename      = "${path.module}/main.zip"
  handler       = "api"
  runtime       = "go1.x"
  role          = join("", aws_iam_role.lambda.*.arn)
  publish       = false

  environment {
    variables = {
      DYNAMODB_TABLE = aws_dynamodb_table.magpie_table.name
    }
  }
  depends_on = [aws_cloudwatch_log_group.api]
}

resource "aws_cloudwatch_log_group" "api" {
  count             = module.this.enabled ? 1 : 0
  name              = "/aws/lambda/${module.api_label.id}"
  retention_in_days = 14
}