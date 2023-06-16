package routine

import (
	"database/sql"
	"time"

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

func createDateTable() {
	conn, err := initDB()
	errPrint(err)

	var now = time.Now()
	addOneDate, _ := time.ParseDuration("24h")
	afterThirty := now.Add(30 * addOneDate)

	for now != afterThirty {
		stmt, err := conn.Prepare("CREATE TABLE IF NOT EXISTS `" + now.Format("2006-01-02") + "` (`ID` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`time` varch(8) NOT NULL,`two` INTEGER NOT NULL,`four` INTEGER NOT NULL,`six` INTEGER NOT NULL,`remark` text);")
		errPrint(err)

		res, err := stmt.Exec()
		_ = res
		errPrint(err)

		createColumn(now.Format("2006-01-02"))

		now = now.Add(addOneDate)
	}
	conn.Close()
}

func createColumn(date string) {
	conn, err := initDB()
	errPrint(err)

	timeArray := [52]string{"09:00", "09:15", "09:30", "09:45",
		"10:00", "10:15", "10:30", "10:45",
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
		res, err := conn.Query("SELECT * FROM `" + date + "`")
		errPrint(err)
		if !res.Next() {
			stmt, err := conn.Prepare("INSERT INTO `" + date + "` (time, two, four, six) VALUES(?,0,0,0)")
			errPrint(err)
			res, err := stmt.Exec(value)
			errPrint(err)
			_ = res
		}
	}
}

func Run() {
	createDateTable()
	// routineTask := cron.New()
	// routineTask.AddFunc("59 23 * * * *", func() {
	// 	createDateTable()
	// })

	// routineTask.Start()
}
