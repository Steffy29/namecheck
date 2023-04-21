package github

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/steffy29/namecheck"
)

var re = regexp.MustCompile("^[-0-9A-Za-z]{3,39}$")

type Github struct {
	Getter namecheck.Getter
}

func (g *Github) IsValid(username string) bool {
	return re.MatchString(username) &&
		!strings.HasPrefix(username, "-") &&
		!strings.HasSuffix(username, "-") &&
		!strings.Contains(username, "--")
}

func (*Github) String() string {
	return "Github"
}

func (g *Github) IsAvailable(username string) (bool, error) {
	endpoint := fmt.Sprintf("https://github.com/%s", username)
	resp, err := g.Getter.Get(endpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return false, nil
	case http.StatusNotFound:
		return true, nil
	default:
		return false, errors.New("unknown availability")
	}
}
