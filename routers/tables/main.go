package tables

import (
	"database/sql"
	"ddic/types"
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

func usedTable() {
	conn, err := initDB()
	errPrint(err)

	now := time.Now()
	dates := now.Format("2006-01-02")
	timeRange, _ := time.ParseDuration("15m")
	nowTime := now.Format("15:04:05")
	afterTime := now.Add(timeRange).Format("15:04:05")

	rows, err := conn.Query("SELECT tableID FROM orderQueue WHERE date=? and time(time) between time(?) and time(?)", dates, nowTime, afterTime)
	errPrint(err)

	for rows.Next() {
		var table int

		err := rows.Scan(&table)
		errPrint(err)

		stmt, err := conn.Prepare("UPDATE tableStatus SET status=1 WHERE tableID=?")
		errPrint(err)

		stmt.Exec(table)
	}

	conn.Close()
}

func releaseTable() {
	conn, err := initDB()
	errPrint(err)

	now := time.Now()
	dates := now.Format("2006-01-02")
	times := now.Format("15:04:05")

	rows, err := conn.Query("SELECT tableID FROM orderQueue WHERE date=? and time(time,'+90 minutes') < time(?) and remark='using'", dates, times)
	errPrint(err)

	for rows.Next() {
		var table int

		err := rows.Scan(&table)
		errPrint(err)

		stmt, err := conn.Prepare("UPDATE tableStatus SET status=0 WHERE tableID=?")
		errPrint(err)

		stmt.Exec(table)
	}
	conn.Close()
}

func cleanTable() {
	conn, err := initDB()
	errPrint(err)

	now := time.Now()
	dates := now.Format("2006-01-02")
	times := now.Format("15:04:05")

	rows, err := conn.Query("SELECT tableID FROM orderQueue WHERE date=? and time(time,'+92 minutes') < time(?) and remark='using'", dates, times)
	errPrint(err)

	for rows.Next() {
		var table int

		err := rows.Scan(&table)
		errPrint(err)

		stmt, err := conn.Prepare("UPDATE tableStatus SET clean=0 WHERE tableID=?")
		errPrint(err)

		stmt.Exec(table)
	}
	conn.Close()
}

func uncleanTable() {
	conn, err := initDB()
	errPrint(err)

	now := time.Now()
	dates := now.Format("2006-01-02")
	times := now.Format("15:04:05")

	rows, err := conn.Query("SELECT tableID FROM orderQueue WHERE date=? and time(time,'+90 minutes') < time(?) and remark='using'", dates, times)
	errPrint(err)

	for rows.Next() {
		var table int

		err := rows.Scan(&table)
		errPrint(err)

		stmt, err := conn.Prepare("UPDATE tableStatus SET status=1 WHERE tableID=?")
		errPrint(err)

		stmt.Exec(table)
	}
	conn.Close()
}

func checkClean() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		rows, err := conn.Query("SELECT tableID FROM tableStatus WHERE clean=0")
		errPrint(err)

		var store = make([]int, 0)

		for rows.Next() {
			var table int
			err = rows.Scan(&table)
			errPrint(err)

			store = append(store, table)
		}

		conn.Close()

		context.JSON(200, gin.H{"CleanTable": store})
	}
}

func getAllStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		conn, err := initDB()
		errPrint(err)

		rows, err := conn.Query("SELECT tableID,status,clean FROM tableStatus")
		errPrint(err)

		var store = make([]types.TableStruct, 0)

		for rows.Next() {
			var table int
			var used int
			var clean int
			err := rows.Scan(&table, &used, &clean)
			errPrint(err)

			store = append(store, types.TableStruct{ID: table, Used: used, Clean: clean})
		}

		context.JSON(200, gin.H{"Status": store})
	}
}

func checkTableUsed() {
	check := cron.New()
	check.AddFunc("*/15 * * * * *", func() {
		usedTable()
		releaseTable()
		uncleanTable()
		cleanTable()
	})
	check.Start()
}

func LoadTableRoutes(e *gin.Engine) {
	checkTableUsed()
	tableRoute := e.Group("/tables")
	{
		tableRoute.GET("/status", getAllStatus())
		tableRoute.GET("/clean", checkClean())
	}
}
