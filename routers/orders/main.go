package orders

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

func emptySeat() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		date := context.Query("date")
		aldult, _ := strconv.Atoi(context.Query("aldult"))
		child, _ := strconv.Atoi(context.Query("child"))
		sumOfPeople := aldult + child

		emptyTime := make([]string, 0)

		fmt.Println(date)
		rows, err := conn.Query("SELECT time,two,four,six FROM `" + date + "`")
		errPrint(err)

		for rows.Next() {
			var getTime string
			var two int
			var four int
			var six int
			rows.Scan(&getTime, &two, &four, &six)
			if sumOfPeople > 36 {
			} else if two < 6 && four < 3 && six < 2 {
				emptyTime = append(emptyTime, getTime)
			} else if sumOfPeople <= 2 && two < 6 {
				emptyTime = append(emptyTime, getTime)
			} else if sumOfPeople <= 4 && (four < 3 || two < 5) {
				emptyTime = append(emptyTime, getTime)
			} else if sumOfPeople <= 6 && ((four < 3 && two < 6) || six < 2) {
				emptyTime = append(emptyTime, getTime)
			} else {
				var needtwo, needfour, needsix = calculateTable(sumOfPeople)
				if two+needtwo <= 4 && four+needfour <= 3 && six+needsix <= 2 {
					emptyTime = append(emptyTime, getTime)
				}
			}
		}

		conn.Close()

		context.JSON(200, gin.H{"Empty": emptyTime})
	}
}

func calculateTable(sumOfPeople int) (int, int, int) {
	var six = sumOfPeople / 6
	sumOfPeople = sumOfPeople % 6

	var four = sumOfPeople / 4
	sumOfPeople = sumOfPeople % 4

	var two = sumOfPeople / 2
	sumOfPeople = sumOfPeople % 4

	_ = sumOfPeople

	return two, four, six
}

func getStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		var date string

		if context.Query("date") != "" {
			date = context.Query("date")
		} else {
			date = time.Now().Format("2006-01-02")
		}

		conn, err := initDB()
		errPrint(err)

		rows, err := conn.Query("SELECT Name,Phone,Gender,Date,Time,Notify FROM orderQueue WHERE date>='" + date + "' AND notify<2")
		errPrint(err)

		customer := make([]types.CustomerStruct, 0)
		for rows.Next() {
			var name string
			var phone string
			var gender int
			var date string
			var time string
			var notify int
			err := rows.Scan(&name, &phone, &gender, &date, &time, &notify)
			errPrint(err)
			customer = append(customer, types.CustomerStruct{Name: name, Phone: phone, Gender: gender, Date: date, Time: time, Notify: notify})
		}
		conn.Close()

		context.JSON(200, gin.H{"Status": customer})
	}
}

func updateNotify() gin.HandlerFunc {
	return func(context *gin.Context) {
		var formCustomer types.PostCustomerStruct
		var customer []types.CustomerStruct
		var conn *sql.DB
		var stmt *sql.Stmt
		var err error

		err = context.ShouldBind(&formCustomer)
		errPrint(err)

		customer = formCustomer.Data

		for i := 0; i < len(customer); i++ {
			conn, err = initDB()
			errPrint(err)

			if customer[i].Notify == 1 {
				stmt, err = conn.Prepare("UPDATE `orderQueue` set notify=1 where name=? and gender=? and phone=? and date=? and time=?;")
			} else {
				stmt, err = conn.Prepare("UPDATE `orderQueue` set notify=2 where name=? and gender=? and phone=? and date=? and time=?;")
			}
			errPrint(err)

			res, err := stmt.Exec(customer[i].Name, customer[i].Gender, customer[i].Phone, customer[i].Date, customer[i].Time)
			errPrint(err)
			_ = res

			conn.Close()
		}

		context.JSON(200, "succes")
	}
}

func orderSeat() gin.HandlerFunc {
	return func(context *gin.Context) {
		var order types.OrderStruct
		var id int64
		var err error

		err = context.BindJSON(&order)
		errPrint(err)

		conn, err := initDB()
		errPrint(err)

		stmt, err := conn.Prepare("INSERT INTO orderQueue(tableID, name, gender, phone, aldult, child, date, time, notify, remark) values(?, ?, ?, ?, ?, ?, ?, ?, 0, ?)")
		errPrint(err)

		for _, value := range order.Table {
			res, err := stmt.Exec(value, order.Name, order.Gender, order.Phone, order.Aldult, order.Child, order.Date, order.Time, order.Remark)
			errPrint(err)

			id, err = res.LastInsertId()
			errPrint(err)

			fmt.Println("Last ID: ", strconv.Itoa(int(id)))
		}
		conn.Close()

		context.JSON(200, gin.H{"序號:": strconv.Itoa(int(id))})
	}
}

func orderFromLive() gin.HandlerFunc {
	return func(context *gin.Context) {
		var order types.CustomerLiveStruct
		var id int64
		var err error

		date := time.Now().Format("2006-01-02")
		hour, minute, second := time.Now().Clock()
		times := strconv.Itoa(hour) + ":" + strconv.Itoa(minute) + ":" + strconv.Itoa(second)

		err = context.BindJSON(&order)
		errPrint(err)

		conn, err := initDB()
		errPrint(err)

		stmt, err := conn.Prepare("INSERT INTO orderQueue(tableID, name, gender, phone, aldult, child, date, time, notify, remark) values(?, ?, ?, ?, ?, ?, ?, ?, 0, ?)")
		errPrint(err)

		for _, value := range order.Table {
			res, err := stmt.Exec(value, order.Name, "2", "0900000000", order.Aldult, order.Child, date, times, "")
			errPrint(err)

			id, err = res.LastInsertId()
			errPrint(err)
		}
		conn.Close()

		context.JSON(200, gin.H{"序號:": strconv.Itoa(int(id))})
	}
}

func LoadOrderRoutes(e *gin.Engine) {

	orderRoute := e.Group("/order")
	{
		orderRoute.GET("/empty", emptySeat())
		orderRoute.GET("/status", getStatus())
		orderRoute.POST("/seat", orderSeat())
		orderRoute.POST("/live", orderFromLive())
		orderRoute.POST("/notify", updateNotify())
	}
}
