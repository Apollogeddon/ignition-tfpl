resource "ignition_alarm_notification_profile" "example" {
  name = "ProductionEmail"
  type = "EmailNotificationProfileType"

  email_config {
    use_smtp_profile = true
    email_profile    = "PrimarySMTP"
  }
}
