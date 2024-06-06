package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandleGetUser(t *testing.T) {
	s := NewServer()
	ts := httptest.NewServer(http.HandlerFunc(s.HandleGetUser))
	reqQty := 1000
	wg := &sync.WaitGroup{}

	for i := 0; i < reqQty; i++ {
		wg.Add(1)
		// Schedules func so it does not wait the func to end to start another iteration
		// Thats why wg is used to sync goroutines with Add, Done and Wait methods
		// Only when all added goroutines are done fmt.Println with the number of db queries runs
		go func(i int) {
			randomId := GenerateRandomID(i)
			if randomId > 100 {
				fmt.Println(randomId)
			}
			url := fmt.Sprintf("%s/?id=%d", ts.URL, randomId)
			res, err := http.Get(url)

			if err != nil {
				t.Error(err)
			}

			user := &User{}
			if err := json.NewDecoder(res.Body).Decode(user); err != nil {
				t.Error(err)
			}
			fmt.Printf("%v\n", user)
			wg.Done()
		}(i)
		// Used to prevent errors from http requests
		time.Sleep(time.Millisecond * 1)
	}

	wg.Wait()
	fmt.Println("Number of db queries: ", s.Db.Queries)
}
