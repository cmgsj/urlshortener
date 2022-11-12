package main

import "urlshortener/pkg/api"

func main() {
	api.NewService().Run()
}
