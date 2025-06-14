package main

import "testing"

func TestIsPrime(t *testing.T) {
	testCases := []struct {
		name     string
		number   int
		expected bool
	}{
		{"negative number", -1, false},
		{"zero", 0, false},
		{"one", 1, false},
		{"two", 2, true},
		{"three", 3, true},
		{"four", 4, false},
		{"large prime", 97, true},
		{"large non-prime", 100, false},
		{"prime number 5", 5, true},
		{"non-prime number 6", 6, false},
		{"prime number 7", 7, true},
		{"non-prime number 9", 9, false},
		{"prime number 13", 13, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isPrime(tc.number); got != tc.expected {
				t.Errorf("isPrime(%d) = %v; want %v", tc.number, got, tc.expected)
			}
		})
	}
}
