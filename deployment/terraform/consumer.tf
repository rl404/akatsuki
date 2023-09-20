resource "kubernetes_deployment" "consumer" {
  metadata {
    name = var.gke_deployment_consumer_name
    labels = {
      app = var.gke_deployment_consumer_name
    }
  }

  spec {
    replicas = 3
    selector {
      match_labels = {
        app = var.gke_deployment_consumer_name
      }
    }
    template {
      metadata {
        labels = {
          app = var.gke_deployment_consumer_name
        }
      }
      spec {
        container {
          name    = var.gke_deployment_consumer_name
          image   = var.gcr_image_name
          command = ["./akatsuki"]
          args    = ["consumer"]
          env {
            name  = "AKATSUKI_CACHE_DIALECT"
            value = var.akatsuki_cache_dialect
          }
          env {
            name  = "AKATSUKI_CACHE_ADDRESS"
            value = var.akatsuki_cache_address
          }
          env {
            name  = "AKATSUKI_CACHE_PASSWORD"
            value = var.akatsuki_cache_password
          }
          env {
            name  = "AKATSUKI_CACHE_TIME"
            value = var.akatsuki_cache_time
          }
          env {
            name  = "AKATSUKI_DB_DIALECT"
            value = var.akatsuki_db_dialect
          }
          env {
            name  = "AKATSUKI_DB_ADDRESS"
            value = var.akatsuki_db_address
          }
          env {
            name  = "AKATSUKI_DB_NAME"
            value = var.akatsuki_db_name
          }
          env {
            name  = "AKATSUKI_DB_USER"
            value = var.akatsuki_db_user
          }
          env {
            name  = "AKATSUKI_DB_PASSWORD"
            value = var.akatsuki_db_password
          }
          env {
            name  = "AKATSUKI_PUBSUB_DIALECT"
            value = var.akatsuki_pubsub_dialect
          }
          env {
            name  = "AKATSUKI_PUBSUB_ADDRESS"
            value = var.akatsuki_pubsub_address
          }
          env {
            name  = "AKATSUKI_PUBSUB_PASSWORD"
            value = var.akatsuki_pubsub_password
          }
          env {
            name  = "AKATSUKI_MAL_CLIENT_ID"
            value = var.akatsuki_mal_client_id
          }
          env {
            name  = "AKATSUKI_CRON_UPDATE_LIMIT"
            value = var.akatsuki_cron_update_limit
          }
          env {
            name  = "AKATSUKI_CRON_FILL_LIMIT"
            value = var.akatsuki_cron_fill_limit
          }
          env {
            name  = "AKATSUKI_CRON_RELEASING_AGE"
            value = var.akatsuki_cron_releasing_age
          }
          env {
            name  = "AKATSUKI_CRON_FINISHED_AGE"
            value = var.akatsuki_cron_finished_age
          }
          env {
            name  = "AKATSUKI_CRON_NOT_YET_AGE"
            value = var.akatsuki_cron_not_yet_age
          }
          env {
            name  = "AKATSUKI_CRON_USER_ANIME_AGE"
            value = var.akatsuki_cron_user_anime_age
          }
          env {
            name  = "AKATSUKI_LOG_JSON"
            value = var.akatsuki_log_json
          }
          env {
            name  = "AKATSUKI_LOG_LEVEL"
            value = var.akatsuki_log_level
          }
          env {
            name  = "AKATSUKI_NEWRELIC_LICENSE_KEY"
            value = var.akatsuki_newrelic_license_key
          }
        }
      }
    }
  }
}
