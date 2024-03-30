resource "aws_sqs_queue" "queue-receiver" {
  name                      = "queue-receiver"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  sqs_managed_sse_enabled   = true

  tags = {
    Name        = "queue-receiver"
    Environment = "DEV"
  }
}