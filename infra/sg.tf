# --- Security Group ---
# ALB

resource "aws_security_group" "alb_sg" {
  name        = "${var.service_name}-alb-sg"
  vpc_id      = aws_vpc.vpc.id
  description = "Security group for ALB"
  tags = {
    Name = "${var.service_name}-alb-sg"
  }
  dynamic "ingress" {
    for_each = [80, 443]
    content {
      protocol    = "tcp"
      from_port   = ingress.value
      to_port     = ingress.value
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Service
resource "aws_security_group" "service_sg" {
  name        = "${var.service_name}-service-sg"
  vpc_id      = aws_vpc.vpc.id
  description = "Security group for service"
  tags = {
    Name = "${var.service_name}-service-sg"
  }
  ingress {
    protocol        = "-1"
    from_port       = 0
    to_port         = 0
    security_groups = [aws_security_group.alb_sg.id]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}
