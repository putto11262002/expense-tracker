terraform {
  required_version = "~> 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

locals {
  web_s3_bucket_name         = "${var.resource_prefix}-web"
  public_subnet_cidr_blocks  = ["10.0.1.0/24"]
  public_subnet_count        = 1
  private_subnet_cidr_blocks = ["10.0.101.0/24", "10.0.102.0/24"]
  private_subnet_count       = 2

}


data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_s3_bucket" "web" {
  bucket = local.web_s3_bucket_name

  tags          = var.common_tags
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "web_public_access_block" {
  bucket = aws_s3_bucket.web.id
}

resource "aws_s3_bucket_policy" "web_allow_public_read" {
  bucket = aws_s3_bucket.web.id
  policy = data.aws_iam_policy_document.web_s3_policy.json
}

data "aws_iam_policy_document" "web_s3_policy" {
  statement {
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.web.arn}/*",
    ]

    principals {
      type        = "*"
      identifiers = ["*"]
    }
  }
}

resource "aws_s3_bucket_website_configuration" "web_s3_website" {
  bucket = aws_s3_bucket.web.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_hostnames = true
  tags                 = var.common_tags
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = aws_vpc.vpc.id
  tags   = var.common_tags
}


resource "aws_subnet" "public_subnet" {
  count  = local.public_subnet_count
  vpc_id = aws_vpc.vpc.id

  cidr_block        = local.public_subnet_cidr_blocks[count.index]
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags              = var.common_tags
}


resource "aws_subnet" "private_subnet" {
  count             = local.private_subnet_count
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = local.private_subnet_cidr_blocks[count.index]
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags              = var.common_tags
}


resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.internet_gateway.id
  }
}

resource "aws_route_table_association" "public_route_table_association" {
  count          = local.public_subnet_count
  route_table_id = aws_route_table.public_route_table.id
  subnet_id      = aws_subnet.public_subnet[count.index].id
}

resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_route_table_association" "private_route_table_association" {
  count          = local.private_subnet_count
  route_table_id = aws_route_table.private_route_table.id
  subnet_id      = aws_subnet.private_subnet[count.index].id
}

resource "aws_security_group" "api_security_group" {
  name        = "${var.resource_prefix}-api-sg"
  description = "Security group for API"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    description = "Allow all traffic through HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow ssh access from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all traffic out"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.common_tags
}

resource "aws_key_pair" "key_pair" {
  key_name   = "${var.resource_prefix}-key-pair"
  public_key = file(var.key_pair_public_key_path)
}

data "aws_ami" "ubuntu-linux" {
  most_recent = true
  owners      = ["099720109477"]
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "api" {
  ami                         = data.aws_ami.ubuntu-linux.id
  instance_type               = var.api_instance_type
  subnet_id                   = aws_subnet.public_subnet[0].id
  key_name                    = aws_key_pair.key_pair.key_name
  vpc_security_group_ids      = [aws_security_group.api_security_group.id]
  tags                        = var.common_tags
  associate_public_ip_address = false
  user_data = templatefile("./api_init.tftpl", {
    db_name = var.database_name
    db_host = aws_db_instance.database.address
    db_port = aws_db_instance.database.port
    db_username = var.database_credentials.username
    db_password = var.database_credentials.password
    allowed_origins = "http://${aws_s3_bucket_website_configuration.web_s3_website.website_endpoint}"
    docker_image = var.api_docker_image
  })
}


resource "aws_security_group" "database_security_group" {
  name        = "${var.resource_prefix}-db-sg"
  description = "Security group for database"
  vpc_id      = aws_vpc.vpc.id
  ingress {
    description     = "Allow MYSQL traffic from API security group"
    from_port       = "3306"
    to_port         = "3306"
    protocol        = "tcp"
    security_groups = [aws_security_group.api_security_group.id]
  }

  tags = var.common_tags
}


resource "aws_db_subnet_group" "database_subnet_group" {
  name       = "${var.resource_prefix}-db-subnet-group"
  subnet_ids = [for subnset in aws_subnet.private_subnet : subnset.id]
  tags       = var.common_tags
}

resource "aws_eip" "api_elastic_ip" {
  depends_on = [aws_internet_gateway.internet_gateway]
  vpc = true
  instance = aws_instance.api.id
}


resource "aws_db_instance" "database" {
  identifier             = "${var.resource_prefix}-database"
  allocated_storage      = 10
  engine                 = "mysql"
  engine_version         = "5.7"
  instance_class         = "db.t2.micro"
  username               = var.database_credentials.username
  password               = var.database_credentials.password
  db_name                = var.database_name
  db_subnet_group_name   = aws_db_subnet_group.database_subnet_group.name
  vpc_security_group_ids = [aws_security_group.database_security_group.id]
  skip_final_snapshot    = true
}