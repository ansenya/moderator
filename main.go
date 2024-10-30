package main

import "moderator/routers"

func main() {
	router := routers.SetupRouter()
	router.Run(":8080")
}
