variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "resource_prefix" {
  description = "Prefix for all resources"
  type        = string
  default     = "shared-expense-tracker"
}

variable "common_tags" {
  description = "Common tags to be applied to all resources"
  type        = map(string)
  default = {
    "Project" = "Shared Expense Tracker"
  }

}


variable "vpc_cidr_block" {
  description = "VIP CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}


variable "api_instance_type" {
  description = "API instance type"
  type        = string
  default     = "t2.micro"

}


variable "database_credentials" {
  description = "Database credentials"
  type        = map(string)
  default = {
    username = "admin"
    password = "password"
  }
  sensitive = true

}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "shared_expense_tracker"
}