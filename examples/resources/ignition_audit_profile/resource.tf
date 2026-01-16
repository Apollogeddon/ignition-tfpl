resource "ignition_audit_profile" "example" {
  name           = "ProductionAudit"
  type           = "database"
  database       = "ProductionDB" # Reference to a database connection
  retention_days = 365
  auto_create    = true
}
