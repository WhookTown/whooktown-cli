module github.com/fredericalix/whooktown-cli

go 1.23

require (
	github.com/fredericalix/whooktown-golang-sdk v0.1.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/fredericalix/whooktown-golang-sdk => ../../sdk/whooktown-golang-sdk
