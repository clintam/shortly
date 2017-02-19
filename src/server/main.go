package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	externalBaseUrl := flag.String("external-url", "http://localhost:8080", "Externally exposed url")
	storage := flag.String("storage", "memory", "Persistent storage. Valid values are memory")
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	flag.Parse()
	s := createServer(*storage, *port, *externalBaseUrl)
	s.Start()
}
