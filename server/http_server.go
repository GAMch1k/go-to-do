package server

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"io"
	"log"
	"net/http"

	"gamch1k.org/todo/database"
	env_manager "gamch1k.org/todo/envmanager"
)

type JSONtasks struct {
	Array []database.Task
}

type ResponseStatus struct {
	Status int `json:"status"`
	Success bool `json:"success"`
	Error_msg string `json:"error_msg"`
}

func api_get_tasks(w http.ResponseWriter, r *http.Request) {
	log.Println("API get tasks request")

	tasks_arr := database.GetTasks(env_manager.GetEnvVariable("DATABASE_PATH"))

	log.Println(tasks_arr)

	json_data, _ := json.Marshal(tasks_arr)

	log.Println("Tasks data converted to JSON")
	io.WriteString(w, string(json_data))
	log.Println("Sent JSON data")
}


func api_post_task(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	var status ResponseStatus

	if text != "" {
		database.InsertTask(env_manager.GetEnvVariable("DATABASE_PATH"), text)
		status = ResponseStatus{
			Status: 200,
			Success: true,
			Error_msg: "",
		}
	} else {
		log.Println("Text parameter is empty")
		status = ResponseStatus{
			Status: 400,
			Success: false,
			Error_msg: "text parameter is empty",
		}
	}
	
	json_data, _ := json.Marshal(status)
	io.WriteString(w, string(json_data))
}


func Start(port string) {
	http.HandleFunc("/api/get_tasks", api_get_tasks)
	http.HandleFunc("/api/post_task", api_post_task)

	log.Println("Starting server on", port)

	err := http.ListenAndServe(port, nil)
	
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatal(err.Error())
		// os.Exit(1)
	}
}