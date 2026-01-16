resource "ignition_database_connection" "example" {
  name        = "production_db"
  type        = "MariaDB"
  connect_url = "jdbc:mariadb://localhost:3306/mydb"
  username    = "dbuser"
  password    = "dbpass"
}
