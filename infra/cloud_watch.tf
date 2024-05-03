# --- Cloud Watch Logs ---
resource "aws_cloudwatch_log_group" "ecs" {
  name              = var.log_path
  retention_in_days = var.log_keep_day
}
