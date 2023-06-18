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
			var table int
			var name string
			var gender int
			var phone string
			var aldult int
			var child int
			var date string
			var time string
			var remark string
			err = rows.Scan(&id, &table, &name, &gender, &phone, &aldult, &child, &date, &time, &remark)
			errPrint(err)

			if len(store) == 0 || !sliceSearch(store, phone, date, time) {
				var tableSlice = make([]int, 0)
				tableSlice = append(tableSlice, table)
				store = append(store, types.OrderStruct{Name: name, Phone: phone, Aldult: aldult, Child: child, Table: tableSlice, Date: date, Time: time, Remark: remark})
			} else {
				sliceInsert(&store, phone, date, time, table)
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
			var table int
			var name string
			var gender int
			var phone string
			var aldult int
			var child int
			var date string
			var time string
			var remark string
			err = rows.Scan(&id, &table, &name, &gender, &phone, &aldult, &child, &date, &time, &remark)
			OrderStatus = append(OrderStatus, types.OrderStatusStruct{TableID: table, Name: name, Gender: gender, Phone: phone, Aldult: aldult, Child: child, Date: date, Time: time, Remark: remark})
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

		stmt, err := conn.Prepare("INSERT INTO orderQueue(tableID, name, gender, phone, aldult, child, date, time, remark) values(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		errPrint(err)

		for _, value := range order.Table {
			res, err := stmt.Exec(value, order.Name, order.Gender, order.Phone, order.Aldult, order.Child, order.Date, order.Time, order.Remark)
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
