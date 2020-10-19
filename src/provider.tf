terraform {
  required_providers {
    cloudflare = {
      version = "2.11.0"
      source  = "cloudflare/cloudflare"
    }

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

provider "cloudflare" {
  api_token = trimspace(file("${path.module}/../.creds/cloudflare_key"))
  account_id = "4d106699f468851a1f005ce8ae96ba5a"
}

provider "google" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}

provider "google-beta" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}
