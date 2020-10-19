terraform {
  required_providers {
    google = {
      version = "3.43.0"
      source  = "hashicorp/google"
    }

    google-beta = {
      version = "3.43.0"
      source  = "hashicorp/google-beta"
    }
  }

  required_version = ">= 0.13"
}

provider "google" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}

provider "google-beta" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}
