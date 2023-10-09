package main

import (
	_ "errors"
	_ "fmt"
	"io"
	"log"
	"os"

	"gamch1k.org/todo/database"
	env_manager "gamch1k.org/todo/envmanager"
	"gamch1k.org/todo/server"
)


func main() {
	log_file, err := os.OpenFile("logs/logs.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer log_file.Close()
	
	mw := io.MultiWriter(os.Stdout, log_file)
	log.SetOutput(mw)

	log.Println("----------- PROGRAM STARTED -----------")

	db_path := env_manager.GetEnvVariable("DATABASE_PATH")

	database.InitDatabase(db_path)

	// database.InsertTask(db_path, "Some test text")

	// fmt.Println(database.GetTasks(db_path))

	server.Start("localhost:" + env_manager.GetEnvVariable("PORT"))

	
}

