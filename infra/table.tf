module "table_label" {
  source     = "cloudposse/label/null"
  version    = "0.25.0"
  context    = module.this.context
  attributes = ["table"]
  enabled    = module.this.enabled
}

resource "aws_dynamodb_table" "magpie_table" {
  name           = module.table_label.id
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "ID"
  range_key      = "SK"

  attribute {
    name = "ID"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }
}