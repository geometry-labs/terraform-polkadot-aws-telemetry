variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  region = var.aws_region
}

variable "public_key_path" {}
variable "private_key_path" {}

resource "random_pet" "this" {}

module "network" {
  source         = "terraform-aws-modules/vpc/aws"
  name           = "telemetry-${random_pet.this.id}"
  cidr           = "10.0.0.0/24"
  azs            = ["${var.aws_region}a"]
  public_subnets = ["10.0.0.0/24"]
  tags = {
    Environment = "CI"
  }
}


module "sg" {
  source = "terraform-aws-modules/security-group/aws"
  name   = "builder"
  vpc_id = module.network.vpc_id
  ingress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
  }]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
  }]
}

module "defaults" {
  source = "../.."
  name   = random_pet.this.id

  private_key_path = var.private_key_path
  public_key_path  = var.public_key_path
  vpc_id           = module.network.vpc_id
  subnet_id        = module.network.public_subnets[0]
}
