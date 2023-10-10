package logs

import (
	"log"
)


func LogError(err error) {
	log.SetPrefix("ERROR ")
	log.Println(err.Error())
	log.SetPrefix("")
}