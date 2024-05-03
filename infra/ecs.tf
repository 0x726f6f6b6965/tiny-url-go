
# --- ECS Cluster ---
resource "aws_ecs_cluster" "ecs_cluster" {
  name = "${var.service_name}-cluster"
}

# --- ECS Task Definition ---
resource "aws_ecs_task_definition" "app_task" {
  family                   = "${var.service_name}-app"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  execution_role_arn       = aws_iam_role.ecs_exec_role.arn
  cpu                      = 256
  memory                   = 512

  container_definitions = jsonencode([
    {
      name         = "${var.service_name}-app"
      image        = "${var.repo_url}:${var.img_tag}"
      essential    = true
      network_mode = "awsvpc"
      environment  = [{ name = "CONFIG", value = file("../deployment/application.yaml") }]
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
        }
      ]
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-region"        = "${var.region}",
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name,
          "awslogs-stream-prefix" = "app"
        }
      },
    }
  ])
}

# --- ECS Service ---
resource "aws_ecs_service" "app_service" {
  name            = var.service_name                     # Name the service
  cluster         = aws_ecs_cluster.ecs_cluster.id       # Reference the created Cluster
  task_definition = aws_ecs_task_definition.app_task.arn # Reference the task that the service will spin up
  launch_type     = "FARGATE"
  desired_count   = var.service_count

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn # Reference the target group
    container_name   = aws_ecs_task_definition.app_task.family
    container_port   = 80 # Specify the container port
  }

  network_configuration {
    subnets          = aws_subnet.public[*].id
    assign_public_ip = true                                                             # Provide the containers with public IPs
    security_groups  = [aws_security_group.service_sg.id, aws_security_group.alb_sg.id] # Set up the security group
  }
  depends_on = [aws_lb_listener.http]
}
