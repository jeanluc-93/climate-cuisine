# resource "aws_sqs_queue" "terraform_queue" {
#   name                      = var.sqs_queue_get_dinner
#   delay_seconds             = 90
#   max_message_size          = 256
#   message_retention_seconds = 86400
#   receive_wait_time_seconds = 10
# }