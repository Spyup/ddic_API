package main

import (
	router "ddic/routes"
	types "ddic/types"

	"github.com/gin-gonic/gin"
)

func createWaitter() []types.WaitterStruct {
	// Create waitter array
	waitter := make([]types.WaitterStruct, 10)
	nameList := [10]string{"JASPER", "EVERLY", "SILAS", "WREN", "GARRETT", "NIA", "MIRANDA", "LEYLA", "NIXON", "JASON"}

	for i := 0; i < 10; i++ {
		waitter[i].Name = nameList[i]
		waitter[i].Distance = 0.0
	}

	return waitter
}

func main() {
	// Create waitter array
	waitterList := createWaitter()

	// Create server
	server := gin.Default()

	// Create route
	server.GET("/", router.GetWaitterList(waitterList))

	// Run server
	server.Run(":9000")
}
