package routes

import (
	types "ddic/types"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   uint64
	Name string
}

func randomWaitterDis(waitterList []types.WaitterStruct) {
	rand.Seed(time.Now().UnixNano())
	var noFormatFloat float64
	for i := 0; i < len(waitterList); i++ {
		noFormatFloat = rand.Float64() * 50
		waitterList[i].Distance, _ = strconv.ParseFloat(strconv.FormatFloat(noFormatFloat, 'f', 2, 64), 64)
	}
}

func GetWaitterList(w []types.WaitterStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		randomWaitterDis(w)
		c.JSON(200, w)
	}
}

func GetWaitterDis(w []types.WaitterStruct) gin.HandlerFunc {
	return func(c *gin.Context) {
		randomWaitterDis(w)
		c.JSON(200, w)
	}

}

func GetWaitterMinDis(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, "")
	}
}
