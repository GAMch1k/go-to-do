package main

import (
	_ "errors"
	"io"
	"log"
	"os"

	"gamch1k.org/todo/database"

	"github.com/joho/godotenv"
)



func get_env_variable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	return os.Getenv(key)
}



func main() {
	log_file, err := os.OpenFile("logs/logs.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer log_file.Close()
	
	mw := io.MultiWriter(os.Stdout, log_file)
	log.SetOutput(mw)

	database.InitDatabase(get_env_variable("DATABASE_PATH"))
}

