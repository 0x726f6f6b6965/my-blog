variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "ec2_ami" {
  type = map(string)
  default = {
    id            = "ami-0f5930086a4d13839"
    instance_type = "t4g.micro"
  }
  description = "Amazon Linux 2"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "vpc_name" {
  description = "Name tag for the VPC"
  default     = "my_blog_vpc"
}

variable "subnet_count" {
  description = "Number of subnets to create"
  default     = 1
}

variable "subnet_cidr" {
  type = list(object({
    cidr_block = string
    az         = string
  }))
  description = "List of subnet CIDR blocks"
  default = [{
    az         = "ap-northeast-1a"
    cidr_block = "10.0.0.0/24"
  }]
}

variable "public_on_launch" {
  type        = bool
  description = "Make the subnet public"
  default     = true
}

variable "key_file" {
  type    = string
  default = "./demo.pub"
}
