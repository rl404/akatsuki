variable "gcp_project_id" {
  type        = string
  description = "GCP project id"
}

variable "gcp_region" {
  type        = string
  description = "GCP project region"
}

variable "gke_cluster_name" {
  type        = string
  description = "GKE cluster name"
}

variable "gke_location" {
  type        = string
  description = "GKE location"
}

variable "gke_pool_name" {
  type        = string
  description = "GKE node pool name"
}

variable "gke_node_preemptible" {
  type        = bool
  description = "GKE node preemptible"
}

variable "gke_node_machine_type" {
  type        = string
  description = "GKE node machine type"
}

variable "gcr_image_name" {
  type        = string
  description = "GCR image name"
}

variable "gke_deployment_consumer_name" {
  type        = string
  description = "GKE deployment consumer name"
}

variable "gke_cron_fill_name" {
  type        = string
  description = "GKE cron fill name"
}

variable "gke_cron_fill_schedule" {
  type        = string
  description = "GKE cron fill schedule"
}

variable "gke_cron_update_name" {
  type        = string
  description = "GKE cron update name"
}

variable "gke_cron_update_schedule" {
  type        = string
  description = "GKE cron update schedule"
}

variable "cloud_run_name" {
  type        = string
  description = "Google cloud run name"
}

variable "cloud_run_location" {
  type        = string
  description = "Google cloud run location"
}

variable "akatsuki_cache_dialect" {
  type        = string
  description = "Cache dialect"
}

variable "akatsuki_cache_address" {
  type        = string
  description = "Cache address"
}

variable "akatsuki_cache_password" {
  type        = string
  description = "Cache password"
}

variable "akatsuki_cache_time" {
  type        = string
  description = "Cache time"
}

variable "akatsuki_db_dialect" {
  type        = string
  description = "Database dialect"
}

variable "akatsuki_db_address" {
  type        = string
  description = "Database address"
}

variable "akatsuki_db_name" {
  type        = string
  description = "Database name"
}

variable "akatsuki_db_user" {
  type        = string
  description = "Database user"
}

variable "akatsuki_db_password" {
  type        = string
  description = "Database password"
}

variable "akatsuki_pubsub_dialect" {
  type        = string
  description = "Pubsub dialect"
}

variable "akatsuki_pubsub_address" {
  type        = string
  description = "Pubsub address"
}

variable "akatsuki_pubsub_password" {
  type        = string
  description = "Pubsub password"
}

variable "akatsuki_mal_client_id" {
  type        = string
  description = "MyAnimeList client id"
}

variable "akatsuki_cron_update_limit" {
  type        = number
  description = "Cron update limit"
}

variable "akatsuki_cron_fill_limit" {
  type        = number
  description = "Cron fill limit"
}

variable "akatsuki_cron_releasing_age" {
  type        = number
  description = "Cron releasing age"
}

variable "akatsuki_cron_finished_age" {
  type        = number
  description = "Cron finished age"
}

variable "akatsuki_cron_not_yet_age" {
  type        = number
  description = "Cron not yet age"
}

variable "akatsuki_cron_user_anime_age" {
  type        = number
  description = "Cron user anime age"
}

variable "akatsuki_log_json" {
  type        = bool
  description = "Log json"
}

variable "akatsuki_log_level" {
  type        = number
  description = "Log level"
}

variable "akatsuki_newrelic_license_key" {
  type        = string
  description = "Newrelic license key"
}
