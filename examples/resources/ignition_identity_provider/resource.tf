resource "ignition_identity_provider" "oidc" {
  name                   = "AzureAD"
  type                   = "oidc"
  client_id              = "my-client-id"
  client_secret          = "my-client-secret"
  provider_id            = "azure-ad"
  authorization_endpoint = "https://login.microsoftonline.com/.../oauth2/v2.0/authorize"
  token_endpoint         = "https://login.microsoftonline.com/.../oauth2/v2.0/token"
  jwk_endpoint           = "https://login.microsoftonline.com/.../discovery/v2.0/keys"
  jwk_endpoint_enabled   = true
}
