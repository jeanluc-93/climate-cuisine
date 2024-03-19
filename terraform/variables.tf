// +-------------------------------------------------------------+
// | Variables saved in Terraform cloud workspace - loaded here. |
// +-------------------------------------------------------------+

variable "default_region" {
    description = "Cape Town region"
    default = "af-south-1"
    type = string
}

// +------------------+
// | Open weather map |
// +------------------+

variable "openWeatherMap_url_key" {
    description = "Reference key to load the Open Weather Map url from Parameter store."
    type = string
}

variable "openWeatherMap_url_value" {
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

// +----------+
// | ClaudeAI |
// +----------+

variable "claude_url_key" {
    description = "Reference key to load the ClaudeAI url-key from Parameter store."
    type = string
}

variable "claude_url_value" {
    description = "URL for the ClaudeAI."
    type = string
}

variable "claude_apiKey_key" {
    description = "Reference key to load the API key-key to access ClaudeAI."
    type = string
}

variable "claude_apiKey_value" {
    description = "The API key to access ClaudeAI."
    type = string
    sensitive = true
}

// +------------+
// | Sqs Queues |
// +------------+

variable "sqs_queue_get_dinner" {
    description = "SQS queue used to push messages to additional lambda to get dinner ideas based off of the received input."
    type = string
}