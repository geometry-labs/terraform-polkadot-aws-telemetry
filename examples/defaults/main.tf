variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  region = var.aws_region
}

variable "public_key_path" {}
variable "private_key_path" {}

resource "random_pet" "this" {}

module "defaults" {
  source = "../.."
  name   = random_pet.this.id

  private_key_path = var.private_key_path
  public_key_path  = var.public_key_path
}
