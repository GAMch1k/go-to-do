package database

import (
	"database/sql"
	"errors"
	"os"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


func InitDatabase(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Println("Creating database file in", path)
		
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("Database created in", path)
	} else {
		log.Printf("Database in %s already exists \n", path)
	}

	sqliteDatabase, _ := sql.Open("sqlite3", path)
	defer sqliteDatabase.Close()
}

