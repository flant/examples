package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

var (
	url = flag.String("u", "", "url")
	n   = flag.Uint("n", 1000, "requests count")
	c   = flag.Uint("c", 100, "concurrent requests")
)

func main() {
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(2)
	}

	concurrentChan := make(chan struct{}, int(*c))
	for i := 0; i < int(*c); i++ {
		concurrentChan <- struct{}{}
	}

	var failedRequests, notOk uint32

	wg := sync.WaitGroup{}
	wg.Add(int(*n))

	start := time.Now()

	var notOkCode uint32

	for i := 0; i < int(*n); i++ {
		<-concurrentChan

		uid := uuid.New().String()
		utn := fmt.Sprintf("%d", time.Now().UnixNano())

		body := strings.NewReader(
			`{
					"id":"` + uid + `",
					"created":` + utn + `,
					"text": ""Test text only for example"
				}`)

		go func() {
			defer func() {
				concurrentChan <- struct{}{}
				wg.Done()
			}()

			resp, err := http.Post(*url, "application/json", body)
			if err != nil {
			    //fmt.Printf("Error :%s\n", err)
                atomic.AddUint32(&failedRequests, 1)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				atomic.AddUint32(&notOk, 1)
				atomic.StoreUint32(&notOkCode, uint32(resp.StatusCode))
				return
			}
		}()

		if i % (int(*n)/10) == 0 {
			fmt.Printf("Completed %d requests\n", i)
		}
	}

	wg.Wait()

	end := time.Now()

	fmt.Printf("\n")
    fmt.Printf("----- Bench results begin -----\n")
	fmt.Printf("Requests per second: %.2f\n", float64(*n)/end.Sub(start).Seconds())
	fmt.Printf("Failed requests: %d\n", failedRequests)
// 	fmt.Printf("Not OK: %d\n", notOk)
// 	fmt.Printf("Not OK code: %d\n", notOkCode)
    fmt.Printf("----- Bench results end -----\n")
}
