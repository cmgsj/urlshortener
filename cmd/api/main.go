package main

import (
	"urlshortener/pkg/services/api"
)

func main() {
	api.NewService().Run()
}
