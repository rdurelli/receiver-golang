resource "aws_s3_bucket" "bucket_1" {
  bucket = "bucket-for-mp3-and-mp4"
  tags = {
    Name        = "Bucket"
    Environment = "Dev"
  }
}