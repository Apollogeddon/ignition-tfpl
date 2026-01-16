resource "ignition_alarm_notification_profile" "example" {
  name = "ProductionEmail"
  type = "email"
  settings = {
    "smtp_profile" = "PrimarySMTP"
  }
}
