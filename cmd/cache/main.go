package main

import "github.com/mike9107/urlshortener/pkg/cache"

func main() {
	cache.NewService().Run()
}
