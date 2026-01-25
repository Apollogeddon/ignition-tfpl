package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func init() {
	// Set default credentials for local testing with Docker if not already set
	if os.Getenv("IGNITION_HOST") == "" {
		_ = os.Setenv("IGNITION_HOST", "http://localhost:8088")
	}
	if os.Getenv("IGNITION_TOKEN") == "" {
		_ = os.Setenv("IGNITION_TOKEN", "terraform:bNxwTt2cyiFUwWFliYY6Fc5flj-AcdqCjfNqn_-Lw8A")
	}
}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ignition": providerserver.NewProtocol6WithError(New("test")()),
}