resource "ignition_gan_settings" "example" {
  require_ssl                    = true
  require_two_way_auth           = false
  allow_incoming                 = true
  security_policy                = "ApprovedOnly"
  websocket_session_idle_timeout = 60000
}
