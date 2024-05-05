package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lock sync.Mutex

func main() {
	url := flag.String("url", "", "endereço a ser testado")
	total := flag.Int("requests", 0, "quantidade total de requisições")
	concurrent := flag.Int("concurrency", 0, "quantidade de requisições concorrentes")

	flag.Parse()

	if *url == "" {
		panic("-url não pode ser vazia")
	}
	if *total == 0 || *concurrent == 0 {
		panic("-requests e -concurrency devem ser maiores que 0")

	}

	var failedRequests = 0
	var rounds = *total / *concurrent
	var lastRound = *total % *concurrent
	var mapOfResponses = make(map[int]int)

	wg := sync.WaitGroup{}
	wg.Add(*total)

	start := time.Now()

	for r := 0; r < rounds; r++ {
		for i := 0; i < *concurrent; i++ {
			go doRequest(*url, &failedRequests, &mapOfResponses, &wg)
		}
	}

	if lastRound > 0 {
		for i := 0; i < lastRound; i++ {
			go doRequest(*url, &failedRequests, &mapOfResponses, &wg)
		}
	}

	wg.Wait()

	executionTime := time.Since(start)

	fmt.Printf("%d requisições executadas com sucesso de um total de %d tentativas.\n", *total-failedRequests, *total)
	fmt.Printf("Tempo de execução: %s\n", executionTime)

	if failedRequests == *total {
		fmt.Println("Todas as requisições falharam.")
	} else {
		fmt.Println("\n\nRelatório detalhado de status code por número de requisições:")
		for key, value := range mapOfResponses {
			fmt.Printf("Código HTTP: %d, Quantidade: %d\n", key, value)
		}
	}
}

func doRequest(url string, failedRequests *int, responses *map[int]int, wg *sync.WaitGroup) {
	lock.Lock()
	defer lock.Unlock()

	req, err := http.Get(url)

	defer wg.Done()

	if err != nil {
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

	defer req.Body.Close()
}
