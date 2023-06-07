package main

import (
	Routers "ddic/routers"
	"fmt"
)

func main() {
	// Create server
	// server := gin.Default()
	router := Routers.InitRouters()

	// Run server
	errorMessage := router.Run(":9000")
	if errorMessage != nil {
		fmt.Println("Service Error")
	}
}
