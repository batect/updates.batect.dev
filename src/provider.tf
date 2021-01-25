// Copyright 2019-2021 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

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
  api_token  = trimspace(file("${path.module}/../.creds/cloudflare_key"))
  account_id = "4d106699f468851a1f005ce8ae96ba5a"
}

provider "google" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}

provider "google-beta" {
  credentials = "${path.module}/../.creds/gcp_service_account_local.json"
}
