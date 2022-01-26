terraform {
  required_version = ">= 0.12"
}

provider "aws" {
  region = "us-west-2"
}

variable "aws_availability_zones" {
  type        = list
  description = "Availability zones to be used"
  default     = ["us-west-2c"]
}

variable "vpc_cidr" {
  description = "The CIDR used for vpc"
  default     = "10.0.0.0/16"
}

variable "tags" {
  description = "Map of tags to add to all resources"
  type        = map
  default     = {}
}

resource "aws_vpc" "konvoy_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = var.tags
}

resource "aws_internet_gateway" "konvoy_gateway" {
  vpc_id = aws_vpc.konvoy_vpc.id

  tags = var.tags
}

resource "aws_subnet" "konvoy_public" {
  vpc_id                  = aws_vpc.konvoy_vpc.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = var.aws_availability_zones[0]

  tags = var.tags
}

resource "aws_route_table" "konvoy_public_rt" {
  vpc_id = aws_vpc.konvoy_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.konvoy_gateway.id
  }

  tags = var.tags
}

resource "aws_route_table_association" "konvoy_public_rta" {
  subnet_id      = aws_subnet.konvoy_public.id
  route_table_id = aws_route_table.konvoy_public_rt.id
}

resource "aws_security_group" "konvoy_ssh" {
  description = "Allow inbound SSH for Konvoy."
  vpc_id      = aws_vpc.konvoy_vpc.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

output "vpc_id" {
  value = aws_vpc.konvoy_vpc.id
}

output "public_subnets" {
  value = aws_subnet.konvoy_public.[0].id
}

output  "security_group_id" {
  value = aws_security_group.konvoy_ssh.[0].id
}
