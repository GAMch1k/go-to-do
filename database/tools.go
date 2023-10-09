package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)


type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Done int `json:"done"`
} 


func InsertTask(path string, text string) {
	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	log.Printf("Inserting %s to the %s (tasks table)", text, path)

	insert_text := `INSERT INTO tasks(text) VALUES (?)`

	query, err := db.Prepare(insert_text)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = query.Exec(text)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Task succesfully inserted")
}


func GetTasks(path string) []Task {
	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	row, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer row.Close()

	var final []Task

	for row.Next() {
		var id int
		var text string
		var done int

		row.Scan(&id, &text, &done)

		final = append(final, Task{
			Id: id,
			Text: text,
			Done: done,
		})
	}

	return final
}


func UpdateTaskDone(path string, id int, status bool) {
	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	log.Printf("Changing task with id %d to %t (tasks table)", id, status)

	query_text := `Update tasks SET done = ? WHERE task_id = ?`

	query, err := db.Prepare(query_text)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = query.Exec(status, id)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Task succesfully updated")
}
