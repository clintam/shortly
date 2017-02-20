package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	serverUrl := flag.String("server", "http://localhost:8080", "The server to use")
	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")
	writeRate := flag.Float64("write-rate", 0.1, "A number between 0 and 1 which represnets the % of operations that are writes")
	iterations := flag.Int("iterations", 1000, "A positive value indicating how many iterations to run (for each client)")
	initialWrites := flag.Int("initial-writes", 100, "A positive value indicating how many initialy shorted urls should be inserted before testing")
	randomSeed := flag.Int64("seed", 42, "A positive value used to seed the random number generator")
	debugMode := flag.Bool("debug", false, "Prints some extra information and opens a HTTP server on port 8081")
	verbose := flag.Bool("verbose", false, "Enables verbose logging")
	flag.Parse()

	rand.Seed(*randomSeed)
	test := MakeTestRun(*serverUrl, *concurrencyLevel, *writeRate, *iterations, *verbose, *initialWrites)

	if *debugMode {
		log.Println("Running in DEBUG mode")
		go func() {
			log.Println(http.ListenAndServe("localhost:8081", nil))
		}()
	}

	test.WarmUp()
	test.Start()
	test.Run()
	test.Finish()
}
