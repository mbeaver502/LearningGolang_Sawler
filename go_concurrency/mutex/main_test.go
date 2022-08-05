package main

import "testing"

// run: $ go test -race .
func Test_updateMessage(t *testing.T) {
	msg = "Hello, world!"

	// test will fail due to data race
	wg.Add(2)
	go updateMessage("Goodbye, cruel world!")
	go updateMessage("xyz")
	wg.Wait()

	if msg != "Goodbye, cruel world!" {
		t.Error("incorrect value in msg")
	}
}
