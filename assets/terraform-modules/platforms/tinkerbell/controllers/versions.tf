# Terraform version and plugin versions

terraform {
  required_version = ">= 0.13"

  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "2.1.2"
    }
    tinkerbell = {
      source  = "tinkerbell/tinkerbell"
      version = "0.1.0"
    }
  }
}
