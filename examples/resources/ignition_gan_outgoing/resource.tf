resource "ignition_gan_outgoing" "example" {
  name    = "site-b"
  host    = "10.20.30.40"
  port    = 8060
  use_ssl = true
}
