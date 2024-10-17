package swagger

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed all:dist
	dist embed.FS
)

func Handler() http.Handler {
	distFS, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(distFS))
}
