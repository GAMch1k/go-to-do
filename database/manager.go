package database

import (
	"database/sql"
	"errors"
	"os"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)


type DB_field struct {
	Name string
	Type string
}

type Database_type struct {
	db *sql.DB
}


func OpenDatabase(path string) (*Database_type) {
	log.Println("Opening database file", path)
	sqliteDatabase, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &Database_type{
		db: sqliteDatabase,
	}
}


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
			log.Printf("Database in %s already exists", path)
	}


	sqliteDatabase := OpenDatabase(path).db
	defer CloseDatabase(sqliteDatabase)
	

	if !CheckIfTableExists(sqliteDatabase, "tasks") {
		fields := []DB_field {
			{
				Name: "task_id",
				Type: "INTEGER PRIMARY KEY AUTOINCREMENT",
			},
			{
				Name: "text",
				Type: "STRING",
			},
			{
				Name: "done",
				Type: "INTEGER DEFAULT 0",
			},
		}
		CreateTable(sqliteDatabase, "tasks", fields)
	}
	
	log.Println("Database fully initialized!")
}


func CloseDatabase(db *sql.DB) {
	db.Close()
	log.Println("Database closed")
}



func CheckIfTableExists(db *sql.DB, name string) bool {
	log.Printf("Checking if table %s exist", name)
	
	query, err := db.Prepare("SELECT name FROM sqlite_master WHERE type='table' AND name=?")
	
	if err != nil {
		log.Fatal(err.Error())
	}
	
	defer query.Close()
	
	var output string
	err = query.QueryRow(name).Scan(&output)
	
	if err == sql.ErrNoRows {
		return false
	}
	
	if err != nil {
		log.Fatal(err.Error())
	}
	
	return true
	
}


func RemoveElementFromFields(slice []DB_field, s int) []DB_field {
	return append(slice[:s], slice[s+1:]...)
}


func CreateTable(db *sql.DB, name string, fields []DB_field) {
	if len(name) <= 0 {
		log.Fatalf("Name of the table %s is empty", name)
	}
	
	if len(fields) <= 0 {
		log.Fatalf("Fields of the table %s is empty", name)
	}
	

	log.Println("Creating new table", name)
	
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s %s)`, name, fields[0].Name, fields[0].Type)
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()

	fields = RemoveElementFromFields(fields, 0)

	log.Println("Created new table:", name)
	log.Println("Creating fields for table", name)
	

	for _, el := range fields {
		createField := fmt.Sprintf("ALTER TABLE %s ADD %s %s", name, el.Name, el.Type)

		statement, err := db.Prepare(createField)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
	}
}
