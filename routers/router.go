package routers

import (
	serviceRouter "ddic/routers/services"
	tableRouter "ddic/routers/tables"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.New()
	serviceRouter.LoadServiceRoutes(router)
	tableRouter.LoadTableRoutes(router)

	return router
}
