package routine

import (
	"database/sql"
	"time"

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

func createDateTable() {

	var now = time.Now()
	addOneDate, _ := time.ParseDuration("24h")
	afterThirty := now.Add(30 * addOneDate)

	for now != afterThirty {
		conn, err := initDB()
		errPrint(err)

		stmt, err := conn.Prepare("CREATE TABLE IF NOT EXISTS `" + now.Format("2006-01-02") + "` (`ID` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`time` varch(8) NOT NULL,`two` INTEGER NOT NULL,`four` INTEGER NOT NULL,`six` INTEGER NOT NULL,`remark` text);")
		errPrint(err)

		res, err := stmt.Exec()
		_ = res
		errPrint(err)

		conn.Close()

		createColumn(now.Format("2006-01-02"))

		now = now.Add(addOneDate)

	}

}

func createColumn(date string) {
	timeArray := [52]string{
		"11:00", "11:15", "11:30", "11:45",
		"12:00", "12:15", "12:30", "12:45",
		"13:00", "13:15", "13:30", "13:45",
		"14:00", "14:15", "14:30", "14:45",
		"15:00", "15:15", "15:30", "15:45",
		"16:00", "16:15", "16:30", "16:45",
		"17:00", "17:15", "17:30", "17:45",
		"18:00", "18:15", "18:30", "18:45",
		"19:00", "19:15", "19:30", "19:45",
		"20:00", "20:15", "20:30", "20:45",
		"21:00", "21:15", "21:30", "21:45"}

	for _, value := range timeArray {
		if value != "" {
			conn, err := initDB()
			errPrint(err)

			res, err := conn.Query("SELECT * FROM `" + date + "` WHERE time='" + value + "';")
			errPrint(err)

			defer res.Close()

			if !res.Next() {
				stmt, err := conn.Prepare("INSERT INTO `" + date + "` (time, two, four, six) VALUES(?,0,0,0)")
				errPrint(err)
				res, err := stmt.Exec(value)
				errPrint(err)
				_ = res
			}
			conn.Close()
		}
	}
}

func Run() {
	routineTask := cron.New()
	routineTask.AddFunc("59 23 * * * *", func() {
		createDateTable()
	})

	routineTask.Start()
}
