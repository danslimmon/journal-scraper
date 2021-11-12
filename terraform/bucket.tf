resource "aws_s3_bucket" "journal_scraper" {
  bucket = var.bucket_name
  acl    = "private"
}
