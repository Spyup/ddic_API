package routers

import (
	customerRouter "ddic/routers/customers"
	orderRouter "ddic/routers/orders"
	serviceRouter "ddic/routers/services"
	tableRouter "ddic/routers/tables"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.New()
	serviceRouter.LoadServiceRoutes(router)
	tableRouter.LoadTableRoutes(router)
	customerRouter.LoadCustomerRoutes(router)
	orderRouter.LoadOrderRoutes(router)

	return router
}
