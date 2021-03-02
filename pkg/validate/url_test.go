package validate

import (
	"fmt"
	"testing"
)

type TestCase struct {
	input    string
	expected bool
}

func TestIsValidURL(t *testing.T) {
	cases := []TestCase{
		{input: "http://google.com", expected: true},
		{input: "google.com", expected: false},
		{input: "www.facebook.com", expected: false},
		{input: "facebook", expected: false},
		{input: "https://facebook.com", expected: true},
	}

	for _, c := range cases {
		_, got := IsValidURL(c.input)
		if got != c.expected {
			t.Log(fmt.Sprintf("should be %+v but got %+v", c.expected, got))
			t.Fail()
		}
	}
}
