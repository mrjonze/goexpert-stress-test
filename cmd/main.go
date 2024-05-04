package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var total = 17
	var concurrent = 10
	var failedRequests = 0
	var rounds = total / concurrent
	var lastRound = total % concurrent
	var mapOfResponses = make(map[int]int)

	if total == 0 || concurrent == 0 {
		panic("Total and Concurrent must be greater than 0")
	}

	var url = "http://google.com"

	wg := sync.WaitGroup{}
	wg.Add(total)

	start := time.Now()

	for r := 0; r < rounds; r++ {
		for i := 0; i < concurrent; i++ {
			go test(i, url, &failedRequests, &mapOfResponses, &wg)
		}
	}

	if lastRound > 0 {
		for i := 0; i < lastRound; i++ {
			go test(i, url, &failedRequests, &mapOfResponses, &wg)
		}

	}

	wg.Wait()

	executionTime := time.Since(start)

	fmt.Printf("%d successful requests from a total of %d attempts.\n", total-failedRequests, total)
	fmt.Printf("Execution time: %s\n", executionTime)

	fmt.Println("Detailed report below:")
	for key, value := range mapOfResponses {
		fmt.Printf("HTTP Code: %d, Amount: %d\n", key, value)
	}

}

func test(i int, url string, failedRequests *int, responses *map[int]int, wg *sync.WaitGroup) {
	req, err := http.Get(url)

	defer wg.Done()

	if err != nil || i == 5 {
		*failedRequests++
		return
	}

	code := req.StatusCode

	value, ok := (*responses)[code]
	if ok {
		(*responses)[code] = value + 1
	} else {
		(*responses)[code] = 1
	}

	fmt.Printf("Request %d -> Response %d\n", i, code)

	defer req.Body.Close()
}
