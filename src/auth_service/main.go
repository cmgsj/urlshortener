package main

import "auth_service/pkg/auth"

func main() {
	auth.NewService().Run()
}
