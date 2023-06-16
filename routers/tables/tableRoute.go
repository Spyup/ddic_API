package tables

import (
	"database/sql"
	"ddic/types"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
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

func createTable() []types.TableStruct {
	// Create waitter array
	table := make([]types.TableStruct, 36)

	for i := 0; i < 36; i++ {
		table[i].ID = i
		// free = 0, using = 1
		table[i].Status = 0
		// clean = 0, using = 1, dirty = 2
		table[i].CleanStatus = 0
	}

	return table
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

func createCRON(table []types.TableStruct) {
	// cron random distance
	randomDis := cron.New()
	randomDis.AddFunc("*/15 * * * * *", func() {
		checkStatus(table)
	})
	randomDis.Start()
}

func checkStatus(table []types.TableStruct) gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		rows, err := conn.Query("SELECT * FROM orderQueue WHERE orderDateTime between ? and ?", "2023-08-01 00:00:00", "2023-12-31 00:00:00")
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

func getStatus(table []types.TableStruct) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, table)
	}
}

func getOrderStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		table := context.Query("tableID")
		currentTime := time.Now()

		rows, err := conn.Query("SELECT * FROM orderQueue WHERE tableID=? AND orderDateTime like ?", table, "%"+currentTime.Format("2006-1-2")+"%")
		errPrint(err)

		OrderStatus := make([]types.OrderStatusStruct, 0)
		for rows.Next() {
			var id int
			var tableID int
			var name string
			var phone string
			var numberOfPeople int
			var dateTime string
			var remark string
			err = rows.Scan(&id, &tableID, &name, &phone, &numberOfPeople, &dateTime, &remark)
			OrderStatus = append(OrderStatus, types.OrderStatusStruct{TableID: tableID, OrderName: name, OrderPhone: phone, NumberOfPeople: numberOfPeople, OrderDateTime: dateTime, Remark: remark})
			errPrint(err)
		}
		conn.Close()

		context.JSON(200, gin.H{"OrderStatus": OrderStatus})
	}
}

func orderSeat() gin.HandlerFunc {
	return func(context *gin.Context) {
		var order types.OrderStruct
		err := context.BindJSON(&order)
		errPrint(err)

		conn, err := initDB()
		errPrint(err)

		stmt, err := conn.Prepare("INSERT INTO orderQueue(tableID, orderName, orderPhone, numberOfPeople, orderDateTime, remark) values(?,?,?,?,datetime(?),?)")
		errPrint(err)

		for _, value := range order.Table {
			res, err := stmt.Exec(strconv.Itoa(value), order.Name, order.Phone, order.NumberOfPeople, order.DateTime, order.Remark)
			errPrint(err)

			id, err := res.LastInsertId()
			errPrint(err)

			fmt.Println("Last ID: ", strconv.Itoa(int(id)))
		}
		conn.Close()
	}
}

func LoadTableRoutes(e *gin.Engine) {
	tableStatus := createTable()
	createCRON(tableStatus)

	tableRoute := e.Group("/tables")
	{
		tableRoute.GET("/status", getStatus(tableStatus))
		tableRoute.GET("/orderStatus", getOrderStatus())
		tableRoute.POST("/order", orderSeat())
		tableRoute.GET("/y", checkStatus(tableStatus))
	}
}
