module github.com/go-obvious/server-example

go 1.23.2

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/render v1.0.3
	github.com/go-obvious/server v0.0.0-00010101000000-000000000000
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.10.0
)

replace github.com/go-obvious/server => ..

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/aws/aws-lambda-go v1.47.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/go-obvious/gateway v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
