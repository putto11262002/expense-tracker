output "web_endpoint" {
  value = aws_s3_bucket_website_configuration.web_s3_website.website_endpoint
}

output "api_endpoint" {
  value = aws_eip.api_elastic_ip.public_dns
}
