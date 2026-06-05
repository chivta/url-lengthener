resource "aws_s3_bucket" "qr" {
  bucket = "${var.project}-${var.environment}-qr"
}

resource "aws_s3_bucket_public_access_block" "qr" {
  bucket = aws_s3_bucket.qr.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
