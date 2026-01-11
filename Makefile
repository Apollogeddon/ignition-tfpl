default: build

build:
	go build -o terraform-provider-ignition.exe

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/apollogeddon/ignition/0.0.1/windows_amd64
	cp terraform-provider-ignition.exe ~/.terraform.d/plugins/registry.terraform.io/apollogeddon/ignition/0.0.1/windows_amd64/terraform-provider-ignition_v0.0.1.exe

test:
	go test ./...
