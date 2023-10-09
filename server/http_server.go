package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"encoding/json"
	_ "fmt"

	"gamch1k.org/todo/database"
	env_manager "gamch1k.org/todo/envmanager"
)

type JSONtasks struct {
	Array []database.Task
}

func api_get_tasks(w http.ResponseWriter, r *http.Request) {
	log.Println("API get tasks request")

	tasks_arr := database.GetTasks(env_manager.GetEnvVariable("DATABASE_PATH"))

	log.Println(tasks_arr)

	json_data, _ := json.Marshal(tasks_arr)

	log.Println(string(json_data))
	io.WriteString(w, string(json_data))
}

func Start(port string) {
	http.HandleFunc("/api/get_tasks", api_get_tasks)

	log.Println("Starting server on", port)

	err := http.ListenAndServe(port, nil)
	
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatal(err.Error())
		// os.Exit(1)
	}
}