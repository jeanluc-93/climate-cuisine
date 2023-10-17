# resource "aws_secretsmanager_secret" "openWeatherMap_secrets" {
#     name = var.openWeatherMap_apiKey_key
#     description = "Secrets' key and value sriven from Terraform."
# }

# resource "aws_secretsmanager_secret_version" "openWeatherMap_apiKey" {
#     secret_id = aws_secretsmanager_secret.openWeatherMap_secrets.id
#     secret_string = var.openWeatherMap_apiKey_value
# }