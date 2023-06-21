package services

import (
	types "ddic/types"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func randomWaitterDis(waitterList []types.WaitterStruct) {
	rand.Seed(time.Now().UnixNano())
	var noFormatFloat float64
	for i := 0; i < len(waitterList); i++ {
		noFormatFloat = rand.Float64() * 50
		waitterList[i].Distance, _ = strconv.ParseFloat(strconv.FormatFloat(noFormatFloat, 'f', 2, 64), 64)
	}
}

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

func createCRON(waitterList []types.WaitterStruct) {
	// cron random distance
	randomDis := cron.New()
	randomDis.AddFunc("*/5 * * * * *", func() {
		randomWaitterDis(waitterList)
	})
	randomDis.Start()
}

func getWaitterMinDis(waitterList []types.WaitterStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		var minDisWaitter types.WaitterStruct
		rand.Seed(time.Now().UnixNano())
		customerTable := rand.Intn(8)

		minDisWaitter.Name = waitterList[0].Name
		minDisWaitter.Distance = waitterList[0].Distance

		for i := 0; i < len(waitterList); i++ {
			if waitterList[i].Distance < minDisWaitter.Distance {
				minDisWaitter.Name = waitterList[i].Name
				minDisWaitter.Distance = waitterList[i].Distance
			}
		}

		response := gin.H{"Customer": customerTable, "Waitter": minDisWaitter}

		c.JSON(200, gin.H{"status": response})
	}
}

func getWaitterList(waitterList []types.WaitterStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		// RandomWaitterDis(waitterList)
		c.JSON(200, gin.H{"WaitterList": waitterList})
	}
}

func LoadServiceRoutes(e *gin.Engine) {
	// Create waitter array
	waitterList := createWaitter()
	randomWaitterDis(waitterList)

	createCRON(waitterList)

	serviceRoute := e.Group("/service")
	{
		serviceRoute.GET("/waitter", getWaitterList(waitterList))
		serviceRoute.GET("/callWaitter", getWaitterMinDis(waitterList))
	}
}
