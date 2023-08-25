package service

import (
	_ "embed"

	"github.com/cmgsj/go-lib/openapi"
)

func init() {
	openapi.Must(openapi.RegisterSchema(schema))
}

var (
	//go:embed urlshortener.swagger.json
	docs []byte

	schema = openapi.Schema{
		Name:        ServiceName,
		ContentJSON: docs,
	}
)
