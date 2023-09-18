resource "aws_ssm_parameter" "openWeatherMapUrl" {
    type = "String"
    name = var.openWeatherMap_apiKey_key
    value = var.openWeatherMap_apiKey_value   
}