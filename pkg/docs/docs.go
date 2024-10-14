package docs

import (
	_ "embed"

	"github.com/cmgsj/go-lib/swagger"

	urlshortenerv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/urlshortener/v1"
)

//go:embed docs.swagger.json
var swaggerDocs []byte

func SwaggerSchema() swagger.Schema {
	return swagger.Schema{
		Name:    urlshortenerv1.URLShortenerService_ServiceDesc.ServiceName,
		Content: swaggerDocs,
	}
}
