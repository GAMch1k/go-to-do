package database

import (
	"log"
	"errors"

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


func GetTaskById(path string, id int) Task {
	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	var task Task

	rows, err := db.Query("SELECT * FROM tasks WHERE task_id = ?", id)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		rows.Scan(&task.Id, &task.Text, &task.Done)
	}

	return task
	
}


func CheckTaskExist(path string, id int) bool {
	res := GetTaskById(path, id)
	
	if res.Id == 0 {
		log.Printf("Task with id %d does not exist", id)
	}

	return res.Id != 0
}


func UpdateTaskDone(path string, id int, status bool) {

	if !CheckTaskExist(path, id) { return }

	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	log.Printf("Changing task with id %d to %t (tasks table)", id, status)

	query_text := `UPDATE tasks SET done = ? WHERE task_id = ?`

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

func DeleteTaskById(path string, id int) (bool, error) {

	if !CheckTaskExist(path, id) { return false, errors.New("Task does not exists") }

	db := OpenDatabase(path).db
	defer CloseDatabase(db)

	log.Println("Deleting task with id", id)

	query_text := `DELETE FROM tasks WHERE task_id = ?`

	query, err := db.Prepare(query_text)
	if err != nil {
		log.SetPrefix("ERROR ")
		log.Println(err.Error())
		log.SetPrefix("")
		return false, err
	}
	_, err = query.Exec(id)
	if err != nil {
		log.SetPrefix("ERROR ")
		log.Println(err.Error())
		log.SetPrefix("")
		return false, err
	}

	log.Println("Task succesfully deleted")
	return true, errors.New("")
}
