
# CORS
# SOURCE: https://github.com/squidfunk/terraform-aws-api-gateway-enable-cors/blob/master/headers.tf
# License-type: MIT

# Copyright (c) 2018-2020 Martin Donath <martin.donath@squidfunk.com>

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to
# deal in the Software without restriction, including without limitation the
# rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
# sell copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
# IN THE SOFTWARE.

locals {
  headers = {
    "Access-Control-Allow-Headers"     = "'${join(",", var.allow_headers)}'"
    "Access-Control-Allow-Methods"     = "'${join(",", var.allow_methods)}'"
    "Access-Control-Allow-Origin"      = "'${var.allow_origin}'"
    "Access-Control-Max-Age"           = "'${var.allow_max_age}'"
    "Access-Control-Allow-Credentials" = "'true'"
  }

  # Pick non-empty header values
  header_values = compact(values(local.headers))

  # Pick names that from non-empty header values
  header_names = matchkeys(
    keys(local.headers),
    values(local.headers),
    local.header_values
  )

  # Parameter names for method and integration responses
  parameter_names = formatlist("method.response.header.%s", local.header_names)

  # Map parameter list to "true" values
  true_list = split("|",
    replace(join("|", local.parameter_names), "/[^|]+/", "true")
  )

  # Integration response parameters
  integration_response_parameters = zipmap(
    local.parameter_names,
    local.header_values
  )

  # Method response parameters
  method_response_parameters = zipmap(
    local.parameter_names,
    local.true_list
  )
}
resource "aws_api_gateway_method" "cors" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "cors" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.cors.http_method
  type        = "MOCK"

  request_templates = {
    "application/json" = "{ \"statusCode\": 200 }"
  }
}

resource "aws_api_gateway_method_response" "cors" {
  depends_on  = [aws_api_gateway_method.cors]
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.cors.http_method
  status_code = 200

  response_parameters = local.method_response_parameters

  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_api_gateway_integration_response" "cors" {
  depends_on          = [aws_api_gateway_integration.cors, aws_api_gateway_method_response.cors]
  rest_api_id         = aws_api_gateway_rest_api.rest_api.id
  resource_id         = aws_api_gateway_resource.proxy.id
  http_method         = aws_api_gateway_method.cors.http_method
  status_code         = 200
  response_parameters = local.integration_response_parameters

  response_templates = {
    "application/json" = "{'message':$context.error.messageString}"
  }
}

resource "aws_api_gateway_gateway_response" "response_4xx" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  response_type = "DEFAULT_4XX"

  response_templates = {
    "application/json" = "{'message':$context.error.messageString}"
  }

}

resource "aws_api_gateway_gateway_response" "response_5xx" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  response_type = "DEFAULT_5XX"

  response_templates = {
    "application/json" = "{'message':$context.error.messageString}"
  }

}