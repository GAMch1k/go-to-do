package main

import (
	_ "errors"
	"os"
	"log"

	"gamch1k.org/todo/database"

	_ "github.com/mattn/go-sqlite3"
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
	database.InitDatabase(get_env_variable("DATABASE_PATH"))
	
}

