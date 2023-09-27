output "web_http_url" {
  value = aws_s3_bucket_website_configuration.web_s3_website.website_endpoint
}

output "api_public_ip" {
  value = aws_instance.api.public_ip
}

output "api_public_dns" {
  value = aws_instance.api.public_dns
}

output "database_endpoint" {
  value = aws_db_instance.database.address
}

output "database_port" {
  value = aws_db_instance.database.port
}
