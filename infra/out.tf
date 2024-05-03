output "app_url" {
  value = aws_alb.alb.dns_name
}
