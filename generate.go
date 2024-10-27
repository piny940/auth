//go:generate go run gen/gorm_gen/main.go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config .oapi-config.yaml spec/schema/@typespec/openapi3/openapi.yaml
//go:generate go run -mod=mod github.com/google/wire/cmd/wire ./...

package main
