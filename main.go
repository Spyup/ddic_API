package main

import (
	Routers "ddic/routers"
	Routines "ddic/routine"
	"fmt"
)

func main() {
	// Create route
	router := Routers.InitRouters()

	// Run routine
	Routines.Run()

	// Run server
	errorMessage := router.Run(":9000")
	if errorMessage != nil {
		fmt.Println("Service Error")
	}
}
