package namecheck

import (
	"fmt"
	"net/http"
)

type Validator interface {
	IsValid(string) bool
}

type Availabler interface {
	IsAvailable(string) (bool, error)
}

type Checker interface {
	Validator
	Availabler
	fmt.Stringer
}

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}
