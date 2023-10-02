output "web_endpoint" {
  value = aws_s3_bucket_website_configuration.web_s3_website.website_endpoint
}

output "api_endpoint" {
  value = aws_lb.api_load_balancer.dns_name
}
