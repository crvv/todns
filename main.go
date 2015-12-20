package main

import (
	"github.com/crvv/todns/dns"
	"github.com/crvv/todns/http"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	logFilename := "log.txt"
	logFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0644)
	if err == nil {
		log.Println("use file log.txt as log output")
		log.SetOutput(logFile)
	}

	http.StartHttpServer()
	dns.Start()
}
