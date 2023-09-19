variable "default_region" {
    description = "Cape Town region"
    default = "af-south-1"
    type = string
}

// Variables saved in Terraform cloud workspace - loaded here.
variable "openWeatherMap_Url_key" {
    description = "Reference key to load the Open Weather Map url from Parameter store."
    type = string
}

variable "openWeatherMap_Url_value" {
    description = "URL for the Open Weather Map api."
    type = string
}

variable "openWeatherMap_apiKey_key" {
    description = "Key that is used to access the correlating value from secrets manager."
    type = string
}

variable "openWeatherMap_apiKey_value" {
    description = "The API Key to access Open Weather Map."
    type = string
    sensitive = true
}

variable "sqs_queue_get_dinner" {
    description = "SQS queue used to push messages to additional lambda to get dinner ideas based off of the received input."
    type = string
}