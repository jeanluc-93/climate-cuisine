resource "aws_iam_role" "lambda_role" {
  name = "LambdaExecutionRole"

  assume_role_policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Sid : "",
        Effect : "Allow",
        Principal : {
          Service : "lambda.amazonaws.com"
        },
        Action : "sts:AssumeRole"
      }
  ] })
}

#IAM Policy to Access DynamoDB in Cape Town Region
resource "aws_iam_policy" "custom_lambda_policy" {
  name        = "custom_lambda_policy"
  description = "AWS IAM Policy for managing AWS lambda role."
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action : [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource : "arn:aws:logs:*:*:*",
        Effect : Allow
      },
      {
        Effect = Allow,
        Action = ["secretsmanager:GetSecretValue"],
        Resource = "arn:aws:secretmanager:${var.default_region}:${local.account_id}:secret:*"
      }
    ]
  })
}

# Policy Attachment on the role.
resource "aws_iam_role_policy_attachment" "attach_lambda_policy_to_iam_role" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.custom_lambda_policy.arn
}

data "archive_file" "zip_go_code" {
  type        = "zip"
  source_dir  = "${path.module}/src/get_weather.go"
  output_path = "${path.module}/build/get_weather.zip"
}

resource "aws_lambda_function" "GetWeather" {
  function_name = "GetWeather"
  description   = "AWS Lambda reaches out to an Open weather map API and gets the weather forecast for the day."
  filename      = "${path.module}/build/get_weather.zip"
  role          = aws_iam_role.lambda_role.arn
  handler       = "main"
  runtime       = "provided.al2"
  memory_size   = 128
  timeout       = 30
  environment {
    variables = {
      REGION = "${var.default_region}",
      SECRET_KEY = "${var.openWeatherMap_apiKey_key}"
    }
  }
  depends_on = [aws_iam_role_policy_attachment.attach_lambda_policy_to_iam_role]
}
