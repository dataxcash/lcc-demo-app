module demo-app

go 1.24.6

require github.com/yourorg/lcc-sdk v0.1.0

require (
	github.com/google/uuid v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/yourorg/lcc-sdk => ../lcc-sdk
