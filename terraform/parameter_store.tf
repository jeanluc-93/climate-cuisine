resource "aws_ssm_parameter" "openWeatherMap_url" {
    type = "String"
    name = var.openWeatherMap_url_key
    value = var.openWeatherMap_url_value
}

resource "aws_ssm_parameter" "openWeatherMap_apiKey" {
    type = "SecureString"
    name = var.openWeatherMap_apiKey_key
    value = var.openWeatherMap_apiKey_value
}

resource "aws_ssm_parameter" "claudeAI_url" {
    type = "String"
    name = var.claude_url_key
    value = var.claude_url_value
}

resource "aws_ssm_parameter" "claudeAi_apiKey" {
    type = "SecureString"
    name = var.claude_apiKey_key
    value = var.claude_apiKey_value
}
