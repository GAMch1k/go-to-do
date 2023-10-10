package database

import (
	"log"
	"errors"

	_ "github.com/mattn/go-sqlite3"

	"gamch1k/todo/logs"
)


type Task struct {
	Id int `json:"id"`
	Text string `json:"text"`
	Done int `json:"done"`
}


func InsertTask(text string) (bool, error) {
	db := OpenDatabase()
	defer CloseDatabase(db)

	log.Printf("Inserting %s to the tasks table", text)

	insert_text := `INSERT INTO tasks(text) VALUES (?)`

	query, err := db.Prepare(insert_text)
	if err != nil {
		logs.LogError(err)
		return false, err
	}
	_, err = query.Exec(text)
	if err != nil {
		logs.LogError(err)
		return false, err
	}

	log.Println("Task succesfully inserted")
	return true, nil
}


func GetTasks() ([]Task, error) {
	db := OpenDatabase()
	defer CloseDatabase(db)

	row, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		logs.LogError(err)
		return []Task{}, err
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

	return final, nil
}


func GetTaskById(id int) (Task, error) {
	db := OpenDatabase()
	defer CloseDatabase(db)

	var task Task

	rows, err := db.Query("SELECT * FROM tasks WHERE task_id = ?", id)

	if err != nil {
		logs.LogError(err)
		return Task{}, err
	}

	for rows.Next() {
		rows.Scan(&task.Id, &task.Text, &task.Done)
	}

	return task, nil
	
}


func CheckTaskExist(id int) (bool, error) {
	res, err := GetTaskById(id)
	
	if res.Id == 0 {
		log.Printf("Task with id %d does not exist", id)
	}

	return res.Id != 0, err
}


func UpdateTaskDone(id int, status bool) (bool, error) {

	exists, _ := CheckTaskExist(id)
	if !exists { return false, nil }

	db := OpenDatabase()
	defer CloseDatabase(db)

	log.Printf("Changing task with id %d to %t (tasks table)", id, status)

	query_text := `UPDATE tasks SET done = ? WHERE task_id = ?`

	query, err := db.Prepare(query_text)
	if err != nil {
		logs.LogError(err)
		return false, err
	}
	_, err = query.Exec(status, id)
	if err != nil {
		logs.LogError(err)
		return false, err
	}

	log.Println("Task succesfully updated")
	return true, nil
}

func DeleteTaskById(id int) (bool, error) {

	exists, _ := CheckTaskExist(id)
	if !exists { return false, errors.New("Task does not exists") }

	db := OpenDatabase()
	defer CloseDatabase(db)

	log.Println("Deleting task with id", id)

	query_text := `DELETE FROM tasks WHERE task_id = ?`

	query, err := db.Prepare(query_text)
	if err != nil {
		logs.LogError(err)
		return false, err
	}
	_, err = query.Exec(id)
	if err != nil {
		logs.LogError(err)
		return false, err
	}

	log.Println("Task succesfully deleted")
	return true, nil
}
