resource "google_storage_bucket" "public" {
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

resource "google_storage_bucket_iam_member" "all_users_viewers" {
  bucket = google_storage_bucket.public.name
  role   = "roles/storage.legacyObjectReader"
  member = "allUsers"
}
