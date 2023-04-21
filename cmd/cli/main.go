package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/steffy29/namecheck"
	"github.com/steffy29/namecheck/github"
	"github.com/steffy29/namecheck/twitter"
)

type Result struct {
	Username  string
	Platform  string
	Valid     bool
	Available bool
	Err       error
}

func check(
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
) {
	defer wg.Done()
	res := Result{
		Username: username,
		Platform: checker.String(),
	}
	res.Valid = checker.IsValid(username)
	// fmt.Printf("validity of %q on %v: %t\n", username, checker, valid)
	if !res.Valid {
		return
	}
	res.Available, res.Err = checker.IsAvailable(username)
	if res.Err != nil {
		fmt.Println(res.Err)
		resultCh <- res
		return
	}
	// fmt.Printf("Availability of %q on %v: %t\n", username, checker, available)
	resultCh <- res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: namecheck <username>")
		os.Exit(1)
	}
	username := os.Args[1]

	tw := twitter.Twitter{Getter: http.DefaultClient}
	gh := github.Github{Getter: http.DefaultClient}

	var checkers []namecheck.Checker
	for i := 0; i < 10; i++ {
		checkers = append(checkers, &tw, &gh)
	}
	resultCh := make(chan Result)

	// username := "ert0ter"
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		go check(checker, username, &wg, resultCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	var results []Result
	for res := range resultCh {
		results = append(results, res)
	}
	fmt.Println(results)
}
