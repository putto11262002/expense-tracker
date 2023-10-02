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
  public_subnet_cidr_blocks  = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnet_count        = 2
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


resource "aws_eip" "ip" {
  vpc      = true
  tags = {
    Name = "t4-elasticIP"
  }
}


resource "aws_nat_gateway" "nat_gateway" {
  allocation_id = "${aws_eip.ip.id}"
  subnet_id     = "${aws_subnet.public_subnet[0].id}"
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


  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_nat_gateway.nat_gateway.id}"
  }
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
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    security_groups = [aws_security_group.api_lb_security_group.id]
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


resource "aws_security_group" "api_lb_security_group" {
  name = "${var.resource_prefix}-api-lb-sg"
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  vpc_id = aws_vpc.vpc.id
}

resource "aws_key_pair" "key_pair" {
  key_name   = "${var.resource_prefix}-key-pair"
  public_key = file(var.key_pair_public_key_path)
}

data "aws_ami" "linux_docker" {
  most_recent = true
  owners      = ["102837901569"]
  filter {
    name   = "name"

    
    values = ["aws-elasticbeanstalk-amzn-*.x86_64-docker-hvm-*"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_launch_configuration" "api_launch_template" {
  name_prefix     = "${var.resource_prefix}-api"
  image_id        = data.aws_ami.linux_docker.id
  instance_type   = var.api_instance_type
  key_name                    = aws_key_pair.key_pair.key_name
  security_groups = [aws_security_group.api_security_group.id]
   
   user_data = templatefile("./api_init.tftpl", {
    db_name = var.database_name
    db_host = aws_db_instance.database.address
    db_port = aws_db_instance.database.port
    db_username = var.database_credentials.username
    db_password = var.database_credentials.password
    allowed_origins = "http://${aws_s3_bucket_website_configuration.web_s3_website.website_endpoint}"
    docker_image = var.api_docker_image
  })
  

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "api_auto_scaling_group" {
  name                 = "${var.resource_prefix}-api-auto-scaling-group"
  min_size             = var.api_autoscale_settings.min
  max_size             = var.api_autoscale_settings.max
  desired_capacity     = var.api_autoscale_settings.desired
  launch_configuration = aws_launch_configuration.api_launch_template.name
  vpc_zone_identifier  = [for subnset in aws_subnet.private_subnet : subnset.id]
  health_check_type    = "ELB"
  target_group_arns = [aws_lb_target_group.api_target_group.arn]
}

resource "aws_lb" "api_load_balancer" {
  name               = "${var.resource_prefix}-api-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.api_lb_security_group.id]
  subnets            = [for subnet in aws_subnet.public_subnet: subnet.id]
}

resource "aws_lb_listener" "api_lb_listener" {
  load_balancer_arn = aws_lb.api_load_balancer.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_target_group.arn
  }
}


 resource "aws_lb_target_group" "api_target_group" {
   name     = "${var.resource_prefix}-api-tg"
   port     = 3000
   protocol = "HTTP"
   vpc_id   = aws_vpc.vpc.id
   health_check {
     enabled = true
     port = 3000
     protocol = "HTTP"
     path = "/api/health-check"
   }
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