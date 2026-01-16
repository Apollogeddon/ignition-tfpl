resource "ignition_smtp_profile" "example" {
  name              = "PrimarySMTP"
  hostname          = "smtp.gmail.com"
  port              = 587
  use_ssl_port      = false
  start_tls_enabled = true
  username          = "alerts@example.com"
  password          = "supersecret"
}
