resource "ignition_store_and_forward" "example" {
  name              = "MainStoreForward"
  time_threshold_ms = 5000
  forward_rate_ms   = 1000
  forwarding_policy = "ForwardAlways"
  data_threshold    = 10000
  batch_size        = 100
  scan_rate_ms      = 500
}
