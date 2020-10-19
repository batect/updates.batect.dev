resource "google_storage_bucket" "terraform_state_bucket" {
  name                        = "${data.google_project.project.name}-public"
  project                     = data.google_project.project.name
  location                    = "us-central1"
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}
