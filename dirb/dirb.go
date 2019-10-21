/*
	Goal:
	- Enumerate web directories

	Features:
	- Scan Efficiency
		- Threads
	- Test Coverage
		- File extension (TODO)
		- Recursive (TODO)
	- Anonymization
		- "Transparent Proxy" (TODO)
		- Wrappable network (TODO)
	- Authentications
		- HTTP Auth
		- Cookie Auth
*/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

const httpTimeOutSeconds = 10

var dictFilename string
var targetURL string
var numThreads int

func init() {
	// Init logging module
	log.SetPrefix("[DIRB] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	// Parse argument flags
	flag.StringVar(&dictFilename, "f", "", "Path to the wordlist")
	flag.StringVar(&targetURL, "u", "", "Target URL")
	flag.IntVar(&numThreads, "t", 1, "Number of concurrent requests per URL")
	flag.Parse()

	if dictFilename == "" {
		flag.PrintDefaults()
		log.Fatal("Dictionary path not defined")
	}
	if targetURL == "" {
		flag.PrintDefaults()
		log.Fatal("Target url not defined")
	}
	if numThreads == 1 {
		log.Println("Default to 1 thread")
	} else {
		if numThreads < 1 || numThreads > 10 {
			log.Fatal("Invalid thread number, must be between 1 to 10")
		}
		log.Printf("Using %d threads\n", numThreads)
	}

	// Validate URL correctness
	u, err := url.Parse(targetURL)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err.Error())
	}

	targetURL = u.String()
	log.Println(targetURL)
}

func main() {
	// Open the dictionary and split words by line
	b, err := ioutil.ReadFile(dictFilename)
	if err != nil {
		panic(err)
	}

	// Use channel as semaphores to limit concurrency
	var sem = make(chan int, numThreads)
	// Use waitgroup to wait for all requests to complete
	var wg sync.WaitGroup

	// Loop through each line and attempt the url
	words := strings.Split(string(b), "\n")
	for _, word := range words {
		// Build URL
		u, _ := url.Parse(targetURL)
		u.Path = path.Join(u.Path, word)

		sem <- 1
		wg.Add(1)

		go func() {
			head(u.String())
			<-sem
			wg.Done()
		}()
	}

	// Wait until the channel is empty
	wg.Wait()
}

func head(targetURL string) {
	var netClient = &http.Client{
		Timeout: time.Second * httpTimeOutSeconds,
	}

	// Check URL availability using HEAD to minimalize response size
	res, err := netClient.Head(targetURL)
	if err != nil {
		log.Println("error accessing url: " + err.Error())
	} else {
		if res.StatusCode == 200 || res.StatusCode == 403 {
			// Found a match
			log.Println("Discovered: " + targetURL)
		}
	}
}
