package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256),
	comment TEXT,
	repeat VARCHAR(128)
);
CREATE INDEX scheduler_date_idx ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
