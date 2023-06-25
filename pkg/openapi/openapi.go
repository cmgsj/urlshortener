package openapi

import (
	"embed"
	"io/fs"
)

var (
	//go:embed docs
	docs embed.FS
)

func Docs() fs.FS {
	return docs
}
