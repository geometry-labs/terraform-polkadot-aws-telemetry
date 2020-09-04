variable "private_key_path" {
  description = "Path to private key"
  type        = string
  default     = ""
}

variable "public_key_path" {
  description = "Path to public key"
  type        = string
  default     = ""
}

variable "subnet_id" {
  description = "ID for the subnet to deploy into"
  type        = string
  default     = null
}

variable "vpc_id" {
  description = "ID for the VPC to deploy into"
  type        = string
  default     = null
}

variable "create_dns" {
  description = "Boolean to enable DNS record creation"
  type        = bool
  default     = false
}

variable "host_name" {
  description = "Hostname for server (i.e. \"telemetry\")"
  type        = string
  default     = "telemetry"
}

variable "domain_name" {
  description = "Root domain name for deployment (as is present in Route53)"
  type        = string
  default     = ""
}

variable "name" {
  description = "A unique id for the deployment"
  type        = string
  default     = ""
}

