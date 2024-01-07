provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source           = "./modules/vpc"
  vpc_cidr         = var.vpc_cidr
  vpc_name         = var.vpc_name
  subnet_count     = var.subnet_count
  subnet_cidr      = var.subnet_cidr
  public_on_launch = var.public_on_launch
}


resource "aws_internet_gateway" "my_blog_gw" {
  vpc_id = module.vpc.vpc_id

  tags = {
    Name = "my_blog_gw"
  }
}

resource "aws_default_route_table" "route_table" {
  default_route_table_id = module.vpc.default_route_table_id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.my_blog_gw.id
  }
  tags = {
    Name = "default route table"
  }
}


resource "aws_security_group" "db_security_group" {
  name        = "PostgreSQL"
  description = "Allow SSH and PostgreSQL inbound traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow_tls"
  }
}

resource "aws_key_pair" "ssh_key" {
  key_name   = "demo-ssh"
  public_key = file(var.key_file)
}

resource "aws_instance" "web" {
  ami           = var.ec2_ami["id"]
  instance_type = var.ec2_ami["instance_type"]
  key_name      = aws_key_pair.ssh_key.key_name
  user_data = templatefile("./install_postgres.sh", {
    pg_hba_file = templatefile("./pg_hba.conf", { allowed_ip = "0.0.0.0/0" }),
  })
  subnet_id                   = module.vpc.subnet_ids[0]
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.db_security_group.id]
  tags = {
    Name = "PostgreSQL"
  }
}


