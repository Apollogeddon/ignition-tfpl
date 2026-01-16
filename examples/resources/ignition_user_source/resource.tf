resource "ignition_user_source" "example" {
  name        = "internal-users"
  type        = "INTERNAL"
  description = "Managed by Terraform"
}
