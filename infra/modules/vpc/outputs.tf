output "vpc_id" {
  value = aws_vpc.my_blog_vpc.id
}

output "subnet_ids" {
  value = aws_subnet.my_blog_subnet[*].id
}

output "default_route_table_id" {
  value = aws_vpc.my_blog_vpc.default_route_table_id
}
