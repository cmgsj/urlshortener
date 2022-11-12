package main

import "urlshortener/pkg/services/urls"

func main() {
	urls.NewService().Run()
}
