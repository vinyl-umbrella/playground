# DynamoDB テーブル
resource "aws_dynamodb_table" "cpu_results" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "S"
  }

  global_secondary_index {
    name            = "TimestampIndex"
    hash_key        = "timestamp"
    projection_type = "ALL"
  }
}

# Lambda 実行ロール
resource "aws_iam_role" "lambda_execution_role" {
  name = "${var.project_name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# CloudWatch Logs 権限
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_execution_role.name
}

# DynamoDB アクセス権限
resource "aws_iam_role_policy" "lambda_dynamodb_policy" {
  name = "dynamodb-policy"
  role = aws_iam_role.lambda_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ]
      Resource = [
        aws_dynamodb_table.cpu_results.arn,
        "${aws_dynamodb_table.cpu_results.arn}/index/*"
      ]
    }]
  })
}

# Lambda 関数を各アーキテクチャとメモリサイズの組み合わせで作成
resource "aws_lambda_function" "func" {
  // var.architectures で foreach したい
  for_each = { for arch in var.architectures : arch => { arch = arch } }

  filename      = each.value.arch == "x86_64" ? var.lambda_zip_x86 : var.lambda_zip_arm
  function_name = "${var.dynamodb_table_name}-${each.value.arch}"
  role          = aws_iam_role.lambda_execution_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = [each.value.arch]
  memory_size   = "128" // tf 外で管理．初期値は 128MB
  timeout       = 30

  lifecycle {
    ignore_changes = [memory_size]
  }

  source_code_hash = filebase64sha256(
    each.value.arch == "x86_64" ? var.lambda_zip_x86 : var.lambda_zip_arm
  )

  environment {
    variables = {
      RESULTS_TABLE_NAME = aws_dynamodb_table.cpu_results.name
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.lambda_basic_execution,
    aws_iam_role_policy.lambda_dynamodb_policy
  ]
}
