resource "ignition_alarm_journal" "example" {
  name         = "MainJournal"
  type         = "DATASOURCE"
  datasource   = "production_db"
  table_name   = "alarm_events"
  min_priority = "Medium"
}
