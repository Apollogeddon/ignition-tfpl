default: build

build:
	go build -o ignition-tfpl.exe

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/apollogeddon/ignition/0.0.1/windows_amd64
	cp ignition-tfpl.exe ~/.terraform.d/plugins/registry.terraform.io/apollogeddon/ignition/0.0.1/windows_amd64/ignition-tfpl_v0.0.1.exe

test:
	go test ./...
