resource "aws_vpc" "my_blog_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = var.vpc_name
  }
}

resource "aws_subnet" "my_blog_subnet" {
  count                   = var.subnet_count
  vpc_id                  = aws_vpc.my_blog_vpc.id
  cidr_block              = var.subnet_cidr[count.index].cidr_block
  availability_zone       = var.subnet_cidr[count.index].az
  map_public_ip_on_launch = var.public_on_launch
  tags = {
    Name = "${var.vpc_name}-subnet-${count.index}"
  }
}
