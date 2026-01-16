resource "ignition_device" "example" {
  name = "Modbus_TCP"
  type = "com.inductiveautomation.ModbusTcpDriver"
  settings = {
    "hostname" = "127.0.0.1"
    "port"     = 502
    "unitId"   = 1
  }
}
