package main

import (
	"testing"
)

func TestHandler(t *testing.T) {

	results, err := lambdaHandler(nil)
	if err != nil {
		t.Errorf("There is an error in the request %s", err)
	}
	if results.RawResponse != nil {
		t.Errorf("\nExpected Result: `Hello from Lambda!`,\nActual Result: `%s`", results)
	}
}
