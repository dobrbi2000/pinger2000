package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {

	url := flag.String("url", "", "URL for sending the request")
	threads := flag.Int("threads", 1, "Number of threads")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout")

	flag.Parse()

	if *url == "" {
		fmt.Println("Please, provide URL")
		return
	}
	fullTime := time.Now()
	errors := make(chan error)
	responses := make(chan time.Duration)

	for i := 0; i < *threads; i++ {
		go func() {
			start := time.Now()
			client := http.Client{
				Timeout: *timeout,
			}
			resp, err := client.Get(*url)

			if err != nil {
				errors <- err
				return
			}
			duration := time.Since(start)
			responses <- duration
			_ = resp.Body.Close()

		}()
	}

	var totalTime time.Duration
	var longest time.Duration
	var shortest time.Duration
	var successReq int

	for i := 0; i < *threads; i++ {
		select {
		case duration := <-responses:
			totalTime += duration
			if duration > longest {
				longest = duration
			}
			if duration < shortest || shortest == 0 {
				shortest = duration
			}
			successReq++
		case err := <-errors:
			fmt.Printf("The request error: %v\n", err)
		}
	}
	totalFullTime := time.Since(fullTime)

	fmt.Printf("The longest request: %v\n", longest)
	fmt.Printf("The shortest request: %v\n", shortest)
	fmt.Printf("The number of successful requests: %v\n", successReq)
	fmt.Printf("Total execution time of all requests: %v\n", totalFullTime)
	fmt.Printf("Average request time: %v\n", totalTime/time.Duration(successReq))

}

/*
test request

 go run main.go -url=https://www.nasa.gov/ -threads=100 -timeout=10s
*/
