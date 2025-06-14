package main

import (
	"reflect"
	"sort"
	"testing"
)

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

func TestFindPrimes_Orchestration(t *testing.T) { // Renamed
	testCases := []struct {
		name       string
		maxNum     int
		numWorkers int
		expected   []int
	}{
		{
			name:       "primes up to 10 with 1 worker",
			maxNum:     10,
			numWorkers: 1,
			expected:   []int{2, 3, 5, 7},
		},
		{
			name:       "primes up to 30 with 4 workers", // New test case
			maxNum:     30,
			numWorkers: 4,
			expected:   []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
		{
			name:       "primes up to 50 with 2 workers", // Another case
			maxNum:     50,
			numWorkers: 2,
			expected:   []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47},
		},
		{
			name:       "primes up to 1 (no primes) with 4 workers",
			maxNum:     1,
			numWorkers: 4,
			expected:   []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Ints(tc.expected) // Ensure expected is sorted

			primes := findPrimes(tc.maxNum, tc.numWorkers)
			// findPrimes already sorts its output

			if !reflect.DeepEqual(primes, tc.expected) {
				t.Errorf("findPrimes(%d, %d) = %v; want %v", tc.maxNum, tc.numWorkers, primes, tc.expected)
			}
		})
	}
}
