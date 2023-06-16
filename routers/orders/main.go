package orders

import (
	"database/sql"
	"fmt"
	"strconv"

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
		numberOfPeople, _ := strconv.Atoi(context.Query("NumberOfPeople"))

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
			if numberOfPeople > 36 {
			} else if two < 6 && four < 3 && six < 2 {
				emptyTime = append(emptyTime, getTime)
			} else if numberOfPeople <= 2 && two < 6 {
				emptyTime = append(emptyTime, getTime)
			} else if numberOfPeople <= 4 && (four < 3 || two < 5) {
				emptyTime = append(emptyTime, getTime)
			} else if numberOfPeople <= 6 && ((four < 3 && two < 6) || six < 2) {
				emptyTime = append(emptyTime, getTime)
			} else {
				var needtwo, needfour, needsix = calculateTable(numberOfPeople)
				if two+needtwo <= 4 && four+needfour <= 3 && six+needsix <= 2 {
					emptyTime = append(emptyTime, getTime)
				}
			}
		}

		conn.Close()

		context.JSON(200, gin.H{"Empty": emptyTime})
	}
}

func calculateTable(numberOfPeople int) (int, int, int) {
	var six = numberOfPeople / 6
	numberOfPeople = numberOfPeople % 6

	var four = numberOfPeople / 4
	numberOfPeople = numberOfPeople % 4

	var two = numberOfPeople / 2
	numberOfPeople = numberOfPeople % 4

	_ = numberOfPeople

	return two, four, six
}

func LoadOrderRoutes(e *gin.Engine) {

	orderRoute := e.Group("/order")
	{
		orderRoute.GET("/empty", emptySeat())
	}
}
