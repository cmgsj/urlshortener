package docs

import (
	_ "embed"

	"github.com/cmgsj/go-lib/openapi"
	"github.com/cmgsj/urlshortener/pkg/service"
)

//go:embed openapi.swagger.json
var docs []byte

func OpenapiSchema() openapi.Schema {
	return openapi.Schema{
		Name:        service.ServiceName,
		ContentJSON: docs,
	}
}
