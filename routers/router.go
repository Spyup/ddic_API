package routers

import (
	customerRouter "ddic/routers/customers"
	serviceRouter "ddic/routers/services"
	tableRouter "ddic/routers/tables"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.New()
	serviceRouter.LoadServiceRoutes(router)
	tableRouter.LoadTableRoutes(router)
	customerRouter.LoadCustomerRoutes(router)

	return router
}
