package main

import (
	"testing"
)

var tests = []struct {
	name     string
	dividend float32
	divisor  float32
	expected float32
	isError  bool
}{
	{
		name:     "valid-data",
		dividend: 100.0,
		divisor:  10.0,
		expected: 10.0,
		isError:  false,
	},
	{
		name:     "invalid-data",
		dividend: 10.0,
		divisor:  0.0,
		expected: 0,
		isError:  true,
	},
	{
		name:     "expect-five",
		dividend: 50.0,
		divisor:  10.0,
		expected: 5.0,
		isError:  false,
	},
	{
		name:     "expect-fraction",
		dividend: 1.0,
		divisor:  4.0,
		expected: 0.25,
		isError:  false,
	},
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		divide(10.0, 3.0)
	}
}

func TestDivide(t *testing.T) {
	for _, test := range tests {
		got, e := divide(test.dividend, test.divisor)

		if test.isError {
			if e == nil {
				t.Error("expected error, did not get")
			}
		} else {
			if e != nil {
				t.Error("did not expect error, got one", e)
			}
		}

		if got != test.expected {
			t.Error("got", got, "want", test.expected)
		}
	}
}
