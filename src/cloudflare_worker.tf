locals {
  cloudflare_zone_id = "b285aeea52df6b888cdee6d2551ebd32"
  # We can't look this up with a data resource without giving access to all zones in the Cloudflare account :sadface:
}

resource "cloudflare_worker_script" "rewrite" {
  name    = "updates_rewrite"
  content = file("worker.js")
}

resource "cloudflare_worker_route" "rewrite" {
  zone_id     = local.cloudflare_zone_id
  pattern     = "updates.batect.dev/*"
  script_name = cloudflare_worker_script.rewrite.name
}

resource "cloudflare_record" "dns" {
  name    = "updates"
  type    = "A"
  zone_id = local.cloudflare_zone_id
  value   = "1.1.1.1" # Dummy value, never used.
  proxied = true
}
