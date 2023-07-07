package main

import (
	Routers "ddic/routers"
	Routines "ddic/routine"
	"fmt"

	"github.com/gin-contrib/cors"
)

func main() {
	// Create route
	router := Routers.InitRouters()

	// Run routine
	Routines.Run()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://shihyan.nmg.cs.thu.edu.tw:8081"},
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type"},
	}))

	// Run server
	errorMessage := router.Run(":9000")
	if errorMessage != nil {
		fmt.Println("Service Error")
	}
}
