package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
	"sync/atomic"
)

// TestRun controls the current state of the test program.
type TestRun struct {
	ServerUrl        string
	StartedAt        time.Time
	ConcurrencyLevel int
	WriteRate        float64
	Iterations       int

	waiting   sync.WaitGroup
	slugToUrl map[string]string
	mutex     sync.RWMutex
	readOps   uint64
	writeOps  uint64
}

// Start starts the test
func (t *TestRun) Start() {
	log.Println("================")
	log.Println(" Starting test ")
	log.Println("================")
	log.Printf("expected server url [%s]", t.ServerUrl)
	log.Printf("concurrency level   [%d]", t.ConcurrencyLevel)
	t.StartedAt = time.Now()
	log.Println("TESTRUN Starting...")
}

// Finish ends the test
func (t *TestRun) Finish() {
	duration := time.Since(t.StartedAt)
	log.Println("================")
	log.Println("All tests passed!")
	log.Println("================")
	log.Printf("TESTRUN finished! (took %dms)", durationInMillis(duration))
	log.Printf("Read operations: %d (%f op/sec)", t.readOps, perSecond(t.readOps, duration))
	log.Printf("Write operations: %d (%f op/sec)", t.writeOps, perSecond(t.writeOps, duration))
	os.Exit(0)
}

func perSecond(count uint64, d time.Duration) float64 {
	sec := float64(d.Nanoseconds()) / float64(time.Second)
	return float64(count) / sec
}

// Fail fails the test
func (t *TestRun) Fail(reason string) {
	duration := time.Since(t.StartedAt)
	log.Println("================")
	log.Println("  Test FAILED!  ")
	log.Println("================")
	log.Printf("Test failed (took %dms)\n%s", durationInMillis(duration), reason)
	os.Exit(0)
	os.Exit(1)
}

//Failf fails the test with a formatted message
func (t *TestRun) Failf(format string, a ...interface{}) {
	t.Fail(fmt.Sprintf(format, a...))
}

//Faile fails the test with the error as its message
func (t *TestRun) Faile(err error) {
	t.Failf("%v", err)
}

//Run executes the test
func (t *TestRun) Run() {
	startedAt := time.Now()

	log.Printf("Spawning %d clients", t.ConcurrencyLevel)
	t.waiting.Add(t.ConcurrencyLevel)
	for i := 0; i < t.ConcurrencyLevel; i++ {
		go t.startClient(i)
	}
	t.waiting.Wait()

	duration := time.Since(startedAt)
	log.Printf("TESTRUN - FINISHED (took %dms %v)", durationInMillis(duration), duration)
}

func (t *TestRun) startClient(num int) {
	client := makeClient(fmt.Sprintf("client %d", num), t)
	for i := 0; i < t.Iterations; i++ {
		randFloat := float64(rand.Intn(100)) / 100.0
		if randFloat <= t.WriteRate {
			t.shorten(client)
		} else {
			t.expand(client)
		}
	}
	t.waiting.Done()

}

func (t *TestRun) shorten(client ShortlyClient) {
	url := randomUrl()
	slug, err := client.Shorten(url)
	if err != nil {
		t.Fail(fmt.Sprintf("Error while shortening [%s]", url, err))
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	atomic.AddUint64(&t.writeOps, 1)
	t.slugToUrl[slug] = url
	log.Printf("shortened %s to %s", url, slug)
}

func (t *TestRun) expand(client ShortlyClient) {
	shortenedCount := len(t.slugToUrl)
	if shortenedCount == 0 {
		return // FIXME?
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	i := rand.Intn(shortenedCount)
	var slug string
	var url string
	for slug, url = range t.slugToUrl {
		if i == 0 {
			break
		}
		i--
	}

	expandedUrl, err := client.Expand(slug)
	if err != nil {
		t.Fail(fmt.Sprintf("Error while expanding [%s]", slug, err))
	}
	if expandedUrl != url {
		t.Fail(fmt.Sprintf("Expected [%s] to expand to [%s] but was [%s]", slug, url, expandedUrl))
	}

	atomic.AddUint64(&t.readOps, 1)
}

//MakeTestRun returns a new instance of a test run.
func MakeTestRun(serverUrl string, concurrencyLevel int, writeRate float64, iterations int) *TestRun {
	return &TestRun{
		ServerUrl:        serverUrl,
		ConcurrencyLevel: concurrencyLevel,
		WriteRate:        writeRate,
		Iterations:       iterations,
		slugToUrl:        make(map[string]string),
	}
}

func randomUrl() string {
	return fmt.Sprintf("http://example.com/%d", rand.Int63())
}

func makeClient(clientName string, t *TestRun) ShortlyClient {
	client := MakeShortlyClient(clientName, t.ServerUrl)

	return client
}

func durationInMillis(d time.Duration) int64 {
	return d.Nanoseconds() / int64(time.Millisecond)
}
