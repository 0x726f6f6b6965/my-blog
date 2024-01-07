variable "vpc_cidr" {
  description = "CIDR block for the VPC"
}

variable "vpc_name" {
  description = "Name tag for the VPC"
}

variable "subnet_count" {
  description = "Number of subnets to create"
}

variable "subnet_cidr" {
  type = list(object({
    cidr_block = string
    az         = string
  }))
  description = "List of subnet CIDR blocks"
}

variable "public_on_launch" {
  type        = bool
  description = "Make the subnet public"
}
