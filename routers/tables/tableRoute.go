package tables

import (
	"database/sql"
	"ddic/types"
	"fmt"
	"strconv"
	"time"

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

func createTable() []types.TableStruct {
	// Create waitter array
	table := make([]types.TableStruct, 36)

	for i := 0; i < 36; i++ {
		table[i].ID = strconv.Itoa(i)
		// free = 0, using = 1
		table[i].Status = 0
		// clean = 0, using = 1, dirty = 2
		table[i].CleanStatus = 0
	}

	return table
}

func getStatus(table []types.TableStruct) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, table)
	}
}

func getOrderStatus(conn *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		table := context.Query("tableID")

		currentTime := time.Now()
		fmt.Println(currentTime.Format("2006-1-2"))

		rows, err := conn.Query("SELECT * FROM orderQueue WHERE tableID=? AND orderDateTime like ?", table, "%"+currentTime.Format("2006-1-2")+"%")
		errPrint(err)

		OrderStatus := make([]types.OrderStruct, 0)
		for rows.Next() {
			var id int
			var tableID string
			var name string
			var phone string
			var numberOfPeople string
			var dateTime string
			var remark string
			err = rows.Scan(&id, &tableID, &name, &phone, &numberOfPeople, &dateTime, &remark)
			OrderStatus = append(OrderStatus, types.OrderStruct{TableID: tableID, OrderName: name, OrderPhone: phone, NumberOfPeople: numberOfPeople, OrderDateTime: dateTime, Remark: remark})

			errPrint(err)
		}
		context.JSON(200, gin.H{"OrderStatus": OrderStatus})
	}
}

func LoadTableRoutes(e *gin.Engine) {
	tableStatus := createTable()
	conn, err := initDB()
	errPrint(err)

	tableRoute := e.Group("/tables")
	{
		tableRoute.GET("/status", getStatus(tableStatus))
		tableRoute.GET("/orderStatus", getOrderStatus(conn))
		// tableRoute.POST("/order")
	}
}
