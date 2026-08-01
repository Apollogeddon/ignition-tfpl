resource "ignition_opc_ua_connection" "example" {
  name            = "LocalOPCUA"
  discovery_url   = "opc.tcp://localhost:4096"
  endpoint_url    = "opc.tcp://localhost:4096"
  security_policy = "None"
  security_mode   = "None"
}
