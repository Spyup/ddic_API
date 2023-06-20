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

func sliceSearch(s []types.OrderStruct, Phone string, Date string, Time string) bool {

	for _, value := range s {
		if value.Phone == Phone && value.Date == Date && value.Time == Time {
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

func sliceInsert(s *[]types.OrderStruct, Phone string, Date string, Time string, Table int) bool {
	for i, value := range *s {
		if value.Phone == Phone && value.Date == Date && value.Time == Time && searchTable(value.Table, Table) {
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
			var table int
			var name string
			var gender int
			var phone string
			var aldult int
			var child int
			var date string
			var time string
			var remark string
			var notify int
			err = rows.Scan(&id, &table, &name, &gender, &phone, &aldult, &child, &date, &time, &remark, &notify)
			errPrint(err)

			if len(store) == 0 || !sliceSearch(store, phone, date, time) {
				var tableSlice = make([]int, 0)
				tableSlice = append(tableSlice, table)
				store = append(store, types.OrderStruct{Name: name, Phone: phone, Gender: gender, Aldult: aldult, Child: child, Table: tableSlice, Date: date, Time: time, Remark: remark, Notify: notify})
			} else {
				sliceInsert(&store, phone, date, time, table)
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
