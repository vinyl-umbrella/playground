data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
data "archive_file" "lambda" {
  type        = "zip"
  source_file = "src/main.py"
  output_path = "src/lambda.zip"
}

resource "aws_lambda_function" "test_lambda" {
  function_name = "example-cold-start"
  handler       = "main.lambda_handler"
  runtime       = "python3.11"
  architectures = ["arm64"]

  role             = aws_iam_role.iam_for_lambda.arn
  filename         = data.archive_file.lambda.output_path
  source_code_hash = data.archive_file.lambda.output_base64sha256


  environment {
    variables = {
      USER = "user"
    }
  }
}

resource "aws_lambda_function_url" "function_url" {
  function_name      = aws_lambda_function.test_lambda.function_name
  authorization_type = "NONE"
}
