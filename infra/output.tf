output "url" {
  value = aws_api_gateway_deployment.rest_api.invoke_url
}

output "ci_id" {
  value = aws_iam_access_key.gh.id
}

output "ci_secret" {
  value = aws_iam_access_key.gh.secret
}

output "ci_arn" {
  value = module.role.arn
}

output "api_function_name" {
  value = join("", aws_lambda_function.api.*.function_name)
}
