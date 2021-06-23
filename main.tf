module "ansible-docker" {
  source = "github.com/insight-infrastructure/terraform-aws-ansible-docker-compose"

  name = var.name

  private_key_path = var.private_key_path
  public_key_path  = var.public_key_path

  additional_roles = ["geometrylabs.substrate_telemetry"]

  create_dns    = var.create_dns
  hostname      = var.host_name
  domain_name   = var.domain_name
  playbook_vars = { "host_fqdn" : "${var.host_name}.${var.domain_name}" }
  open_ports    = [80, 8000]
  subnet_id     = var.subnet_id
  vpc_id        = var.vpc_id
}
