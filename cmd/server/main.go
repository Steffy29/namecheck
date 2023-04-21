package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/steffy29/namecheck"
	"github.com/steffy29/namecheck/github"
)

type Result struct {
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	Valid     bool   `json:"is_valid"`
	Available bool   `json:"is_available"`
}

type User struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

type Stats struct {
	Username string `json:"username"`
	Count    int    `json:count`
}

var m = make(map[string]uint)
var mu sync.Mutex

func check(
	ctx context.Context,
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
	errorCh chan<- error,
) {
	defer wg.Done()
	res := Result{
		Username: username,
		Platform: checker.String(),
	}
	res.Valid = checker.IsValid(username)
	if !res.Valid {
		select {
		case <-ctx.Done():
			return
		case resultCh <- res:
		}
		return
	}
	available, err := checker.IsAvailable(username)
	if err != nil {
		fmt.Println(err)
		select {
		case <-ctx.Done():
			return
		case errorCh <- err:
		}
		return
	}
	res.Available = available
	select {
	case <-ctx.Done():
		return
	case resultCh <- res:
	}

}

func handleHello(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello world!\n")
}

func handleCheck(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	username := req.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad request")
		return
	}
	mu.Lock()
	m[username]++
	mu.Unlock()

	// tw := twitter.Twitter{Getter: http.DefaultClient}
	gh := github.Github{Getter: http.DefaultClient}

	var checkers []namecheck.Checker
	for i := 0; i < 10; i++ {
		checkers = append(checkers, &gh)
	}
	resultCh := make(chan Result)
	errorCh := make(chan error)

	// username := "ert0ter"
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	for _, checker := range checkers {
		wg.Add(1)
		go check(ctx, checker, username, &wg, resultCh, errorCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	var results []Result
	var finished bool
	for !finished {
		select {
		case err := <-errorCh:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal server error\n", err)
			cancel()
			return
		case res, ok := <-resultCh:
			if !ok {
				finished = true
				continue
			}
			results = append(results, res)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal server error")
	}

}

func handleStats(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	mu.Lock()
	if err := enc.Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal server error")
	}
	mu.Unlock()
}

func main() {
	router := httprouter.New()
	router.GET("/hello", handleHello)
	router.GET("/check", handleCheck)
	router.GET("/stats", handleStats)

	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
