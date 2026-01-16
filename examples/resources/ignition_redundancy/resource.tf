resource "ignition_redundancy" "example" {
  role                 = "Master"
  active_history_level = "Full"
  join_wait_time       = 60
  recovery_mode        = "Automatic"
  
  gateway_network_setup = {
    host       = "backup-gateway"
    port       = 8060
    enable_ssl = true
  }
}
