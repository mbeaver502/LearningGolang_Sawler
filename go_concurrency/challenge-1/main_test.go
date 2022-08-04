package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "Hello, world!"
	printMessage()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, msg) {
		t.Error("expected to find", msg)
	}
}

func Test_updateMessage(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go updateMessage("test", &wg)
	wg.Wait()

	if msg != "test" {
		t.Errorf("got %s; want %s", msg, "test")
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := strings.Split(string(result), "\n")

	if output[0] != "Hello, universe!" && output[1] != "Hello, cosmos!" && output[2] != "Hello, world!" {
		t.Error("invalid output", output)
	}

	os.Stdout = stdOut
}
