package server

import (
	"encoding/json"
	"errors"
	_ "fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"gamch1k/todo/database"
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

	tasks_arr, err := database.GetTasks()

	if err != nil {
		status := ResponseStatus{
			Status: 400,
			Success: false,
			Error_msg: "something wrong with getting tasks: " + err.Error(),
		}

		json_data, _ := json.Marshal(status)
		io.WriteString(w, string(json_data))
		return
	}

	json_data, _ := json.Marshal(tasks_arr)

	log.Println("Tasks data converted to JSON")
	io.WriteString(w, string(json_data))
	log.Println("Sent JSON data")
}


func api_post_task(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")

	var status ResponseStatus

	if text != "" {
		res, err := database.InsertTask(text)
		if res {
			status = ResponseStatus{
				Status: 200,
				Success: true,
				Error_msg: "",
			}
		} else {
			status = ResponseStatus{
				Status: 400,
				Success: false,
				Error_msg: "something went wrong: " + err.Error(),
			}
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


func api_update_task(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	complete := r.URL.Query().Get("status")

	var status ResponseStatus

	if id != "" && complete != "" {
		id_int, err := strconv.Atoi(id)

		if err != nil {
			status := ResponseStatus{
				Status: 400,
				Success: false,
				Error_msg: "something wrong with converting id field: " + err.Error(),
			}
			json_data, _ := json.Marshal(status)
			io.WriteString(w, string(json_data))
			return
		}

		var status_bool bool
		if complete == "true" {
			status_bool = true
		} else {
			status_bool = false
		}

		res, err := database.UpdateTaskDone(id_int, status_bool)

		if res {
			status = ResponseStatus{
				Status: 200,
				Success: true,
				Error_msg: "",
			}
		} else {
			status = ResponseStatus{
				Status: 400,
				Success: false,
				Error_msg: "something went wrong: " + err.Error(),
			}
		}
		
	} else {
		log.Println("Some parameters are missing")
		status = ResponseStatus{
			Status: 400,
			Success: false,
			Error_msg: "some parameters are missing",
		}
	}
	
	json_data, _ := json.Marshal(status)
	io.WriteString(w, string(json_data))
}


func api_delete_task(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var status ResponseStatus

	if id != "" {
		id_int, err := strconv.Atoi(id)
		if err != nil {
			status := ResponseStatus{
				Status: 400,
				Success: false,
				Error_msg: "something wrong with converting id field: " + err.Error(),
			}
			json_data, _ := json.Marshal(status)
			io.WriteString(w, string(json_data))
			return
		}
		res, err := database.DeleteTaskById(id_int)
		
		if res {
			status = ResponseStatus{
				Status: 200,
				Success: true,
				Error_msg: "",
			}
		} else {
			status = ResponseStatus{
				Status: 400,
				Success: false,
				Error_msg: "something went wrong: " + err.Error(),
			}
		}

	} else {
		log.Println("Id parameter is empty")

		status = ResponseStatus{
			Status: 400,
			Success: false,
			Error_msg: "id parameter is empty",
		}
	}

	json_data, _ := json.Marshal(status)
	io.WriteString(w, string(json_data))
}


func Start(port string) {
	http.HandleFunc("/api/get_tasks", api_get_tasks)		// No parameters needed
	http.HandleFunc("/api/post_task", api_post_task)		// text={text of the task} parameter
	http.HandleFunc("/api/update_task", api_update_task)	// id={task_id}&status={bool} parameters 
	http.HandleFunc("/api/delete_task", api_delete_task)	// id={task_id} parameter

	log.Println("Starting server on", port)

	err := http.ListenAndServe(port, nil)
	
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatal(err.Error())
		// os.Exit(1)
	}
}
