terraform {
  backend "s3" {
    bucket  = "my-terraform-tfstate-ap-northeast-1"
    region  = "ap-northeast-1"
    key     = "go-lamdba-tutorial.tfstate"
    encrypt = true
  }
}