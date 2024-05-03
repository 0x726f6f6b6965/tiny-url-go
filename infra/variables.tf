variable "repo_url" {
  type    = string
  default = "hello-world"
}
variable "img_tag" {
  type    = string
  default = "latest"
}

variable "region" {
  type    = string
  default = "ap-northeast-1"
}

variable "az_count" {
  type    = number
  default = 2
}

variable "log_path" {
  type    = string
  default = "/ecs/tiny-url"
}

variable "log_keep_day" {
  type    = number
  default = 7
}

variable "service_count" {
  type    = number
  default = 1
}

variable "service_name" {
  type    = string
  default = "tiny-url"
}
