package twitter_test

import (
	"testing"

	"github.com/steffy29/namecheck/twitter"
)

func TestIsValid(t *testing.T) {
	type TestCase struct {
		desc     string
		username string
		want     bool
	}
	testCases := []TestCase{
		{"contains 'Twitter'", "jub0bsOnTwitter", false},
		{"too short", "foo", false},
		{"too long", "dfdffsfdsfdsfdsfdsfdsfdfdsfdsfdsfdsfdsfdsfdsfdsfdsdfdfdsfdsfdsfdsfdsfdfdsfdsf", false},
		{"contains illegal chars", "f^^oo", false},
		{"all good", "ert0ter", true},
	}
	for _, tc := range testCases {
		tw := twitter.Twitter{}
		got := tw.IsValid(tc.username)
		if got != tc.want {
			t.Errorf("twitter.IsValid(%q): got %t; want %t", tc.username, got, tc.want)
		}
	}
}
