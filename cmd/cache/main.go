package main

import "urlshortener/pkg/services/cache"

func main() {
	cache.NewService().Run()
}
