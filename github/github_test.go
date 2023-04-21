package github_test

import (
	"testing"

	"github.com/steffy29/namecheck/github"
)

func TestUsernameTooShort(t *testing.T) {
	const (
		username = "ab"
		want     = false
	)

	gh := github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("github.IsValid(%q): got %t; want %t", username, got, want)
	}
}

func TestUsernameTooLong(t *testing.T) {
	const (
		username = "skjhasdhfghsdhfgkjdhsfkgjhdsfkskdjfhskjdfhskhdfksjdhfsjkdhfhkdfg"
		want     = false
	)
	gh := github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("github.IsValid(%q): got %t; want %t", username, got, want)
	}
}

func TestUsernameIllegalChars(t *testing.T) {
	const (
		username = "fo^^o"
		want     = false
	)
	gh := github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("github.IsValid(%q): got %t; want %t", username, got, want)
	}
}

func TestUsernameAllGood(t *testing.T) {
	const (
		username = "jub0bs"
		want     = true
	)
	gh := github.Github{}
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("github.IsValid(%q): got %t; want %t", username, got, want)
	}
}
