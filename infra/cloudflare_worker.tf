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

locals {
  # We can't look this up with a data resource without giving access to all zones in the Cloudflare account :sadface:
  cloudflare_zone_id = "b285aeea52df6b888cdee6d2551ebd32"
  domain_name        = "${var.subdomain}.${var.root_domain}"
}

resource "cloudflare_worker_script" "rewrite" {
  name    = "updates_rewrite-${replace(local.domain_name, ".", "_")}"
  content = file("worker.js")
}

resource "cloudflare_worker_route" "rewrite" {
  zone_id     = local.cloudflare_zone_id
  pattern     = "${local.domain_name}/*"
  script_name = cloudflare_worker_script.rewrite.name
}

resource "cloudflare_record" "dns" {
  name    = var.subdomain
  type    = "A"
  zone_id = local.cloudflare_zone_id
  value   = "192.0.2.1" # Dummy value, never used.
  proxied = true
}
