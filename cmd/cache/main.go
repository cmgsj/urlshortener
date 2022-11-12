package main

import "urlshortener/pkg/cache"

func main() {
	cache.NewService().Run()
}
