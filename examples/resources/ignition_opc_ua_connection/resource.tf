resource "ignition_opc_ua_connection" "example" {
  name = "LocalOPCUA"
  type = "com.inductiveautomation.OpcUaServerType"
  endpoint = {
    discovery_url   = "opc.tcp://localhost:4096"
    endpoint_url    = "opc.tcp://localhost:4096"
    security_policy = "http://opcfoundation.org/UA/SecurityPolicy#None"
    security_mode   = "None"
  }
}
