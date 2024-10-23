terraform {
  required_version = "= 1.9.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.69.0"
    }
  }
  backend "local" {
    path = ".terraform/terraform.tfstate"
  }
}

provider "aws" {
  region                      = "ap-northeast-1"
  access_key                  = "oteleport0000"
  secret_key                  = "oteleport0000"
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  s3_use_path_style           = true
  endpoints {
    s3 = "http://localhost:9000"
  }
}

resource "aws_s3_bucket" "main" {
  bucket = "oteleport-local"
  timeouts {
    create = "30s"
    update = "30s"
    delete = "30s"
  }
}
