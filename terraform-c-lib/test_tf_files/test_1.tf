terraform {
  backend "gcs" {
    bucket = "tf-state-bucket-name"

    prefix = "terraform/state"
  }
}

resource "google_service_account" "workstation_service_account" {
  account_id   = "workstation-service-account"
  display_name = "Workstation Service Account"
  project      = "my-google-project"
}

resource "google_compute_instance_template" "workstation_template" {
  name        = "workstation-template"
  description = "This template is used to create workstation instances."
  project     = "my-google-project"

  tags = ["foo", "bar"]

  labels = {
    environment = "dev"
  }

  instance_description = "description assigned to instances"
  machine_type         = "n1-standard-8"
  can_ip_forward       = false

  scheduling {
    automatic_restart   = false
    on_host_maintenance = "TERMINATE"
  }

  guest_accelerator {
    type  = "nvidia-tesla-p4"
    count = 1
  }


  // Create a new boot disk from an image
  disk {
    source_image = "debian-cloud/debian-10"
    auto_delete  = true
    boot         = true
    disk_size_gb = 100
    // backup the disk every day
    //resource_policies = [google_compute_resource_policy.daily_backup.id]
  }

  network_interface {
    network = "default"
    // This adds an empemeral external IP
    access_config {}

  }

  metadata = {
    foo = "bar"
  }

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.workstation_service_account.email
    scopes = ["cloud-platform"]
  }
}

/*resource "google_compute_resource_policy" "daily_backup" {
  name   = "every-day-4am"
  region = "us-central1"
  project = "hallowed-forge-300500"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "04:00"
      }
    }
  }
}*/

resource "google_compute_instance_from_template" "workstation" {
  name    = "workstation-one"
  zone    = "us-central1-a"
  project = "my-google-project"

  source_instance_template = google_compute_instance_template.workstation_template.id

}



