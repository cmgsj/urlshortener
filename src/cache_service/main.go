package main

import "cache_service/pkg/cache"

func main() {
	cache.NewService().Run()
}
