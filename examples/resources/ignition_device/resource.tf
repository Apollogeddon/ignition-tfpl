resource "ignition_device" "example" {
  name = "Simulator1"
  type = "ProgrammableSimulatorDevice"
  parameters = jsonencode({
    repeat           = false
    legacyMode       = false
    timeIntervalRate = 1000
  })
}
