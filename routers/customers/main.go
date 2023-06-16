package customers

import (
	"database/sql"
	"ddic/types"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func errPrint(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./ddic.db")
	errPrint(err)

	return db, err
}

func sliceSearch(s []types.OrderStruct, Phone string, DateTime string) bool {

	for _, value := range s {
		if value.Phone == Phone && value.DateTime == DateTime {
			return true
		}
	}
	return false
}

func searchTable(s []int, ID int) bool {
	for _, value := range s {
		if value == ID {
			return false
		}
	}
	return true
}

func sliceInsert(s *[]types.OrderStruct, Phone string, DateTime string, Table int) bool {
	for i, value := range *s {
		if value.Phone == Phone && value.DateTime == DateTime && searchTable(value.Table, Table) {
			value.Table = append(value.Table, Table)
			(*s)[i].Table = value.Table
			return true
		}
	}
	return false
}

func checkStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		rows, err := conn.Query("SELECT * FROM orderQueue")
		errPrint(err)

		var store = make([]types.OrderStruct, 0)

		for rows.Next() {
			var id int
			var tableID int
			var name string
			var phone string
			var numberOfPeople int
			var dateTime string
			var remark string
			err = rows.Scan(&id, &tableID, &name, &phone, &numberOfPeople, &dateTime, &remark)
			errPrint(err)

			if len(store) == 0 || !sliceSearch(store, phone, dateTime) {
				var tableSlice = make([]int, 0)
				tableSlice = append(tableSlice, tableID)
				store = append(store, types.OrderStruct{Name: name, Phone: phone, NumberOfPeople: numberOfPeople, Table: tableSlice, DateTime: dateTime, Remark: remark})
			} else {
				sliceInsert(&store, phone, dateTime, tableID)
			}

		}
		conn.Close()
		context.JSON(200, gin.H{"OrderStatus": store})
	}
}

func LoadCustomerRoutes(e *gin.Engine) {

	customerRoute := e.Group("/customer")
	{
		customerRoute.GET("/status", checkStatus())
	}
}
