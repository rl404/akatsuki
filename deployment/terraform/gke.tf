resource "google_container_cluster" "cluster" {
  name                     = var.gke_cluster_name
  location                 = var.gke_location
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "pool" {
  name       = var.gke_pool_name
  cluster    = google_container_cluster.cluster.id
  node_count = 1

  node_config {
    spot         = var.gke_node_preemptible
    machine_type = var.gke_node_machine_type
  }
}
