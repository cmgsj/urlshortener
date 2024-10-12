package docs

import (
	_ "embed"

	"github.com/cmgsj/go-lib/swagger"

	"github.com/cmgsj/urlshortener/pkg/service"
)

//go:embed docs.swagger.json
var swaggerDocs []byte

func SwaggerSchema() swagger.Schema {
	return swagger.Schema{
		Name:    service.ServiceName,
		Content: swaggerDocs,
	}
}
