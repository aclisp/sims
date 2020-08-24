locals {
  config_name = "go.micro.config"
  config_port = 8080
  config_labels = merge(
    local.common_labels,
    {
      "name" = local.config_name
    }
  )
  config_annotations = merge(
    local.common_annotations,
    {
      "name" = local.config_name
    }
  )
  config_env = merge(
    local.common_env_vars,
    {
      "MICRO_SERVER_ADDRESS" = "0.0.0.0:8080"
      "MICRO_AUTH"           = "jwt"
    }
  )
}

module "config_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.config_name
}

resource "kubernetes_secret" "config_cert" {
  metadata {
    name        = "${replace(local.config_name, ".", "-")}-cert"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.config_labels
    annotations = local.config_annotations
  }
  data = {
    "cert.pem" = module.config_cert.cert_pem
    "key.pem"  = module.config_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "config" {
  metadata {
    name        = replace(local.config_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.config_labels
    annotations = local.config_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.config_labels
    }
    template {
      metadata {
        labels = local.config_labels
      }
      spec {
        container {
          name = replace(local.config_name, ".", "-")
          dynamic "env" {
            for_each = local.config_env
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name = "MICRO_AUTH_PUBLIC_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.micro_keypair.metadata[0].name
                key  = "public"
              }
            }
          }
          env {
            name = "MICRO_AUTH_PRIVATE_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.micro_keypair.metadata[0].name
                key  = "private"
              }
            }
          }
          args              = ["config"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.config_port
            name           = "config-port"
          }
          volume_mount {
            mount_path = "/etc/micro/certs"
            name       = "certs"
          }
          volume_mount {
            mount_path = "/etc/micro/ca"
            name       = "platform-ca"
          }
        }
        volume {
          name = "platform-ca"
          secret {
            secret_name  = kubernetes_secret.platform_ca.metadata[0].name
            default_mode = "0600"
            items {
              key  = "ca.pem"
              path = "ca.pem"
            }
          }
        }
        volume {
          name = "certs"
          secret {
            default_mode = "0600"
            secret_name  = kubernetes_secret.config_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "config" {
  metadata {
    name        = replace(local.config_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.config_labels
    annotations = local.config_annotations
  }
  spec {
    port {
      port        = local.config_port
      target_port = local.config_port
    }
    selector = local.config_labels
  }
}
