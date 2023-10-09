package server

import (
	"errors"
	"io"
	"log"
	"net/http"
)

func api_get_tasks(w http.ResponseWriter, r *http.Request) {
	log.Println("API get tasks request")
	io.WriteString(w, "Get tasks")
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