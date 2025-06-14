package main

import (
	"reflect"
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

func TestSieveOfEratosthenesSequential(t *testing.T) {
	testCases := []struct {
		name     string
		maxNum   int
		expected []int
	}{
		{"primes up to 10", 10, []int{2, 3, 5, 7}},
		{"primes up to 20", 20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{"primes up to 2", 2, []int{2}},
		{"primes up to 1", 1, []int{}},
		{"primes up to 0", 0, []int{}},
		{"primes up to 30", 30, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// `expected` is already sorted. `sieveOfEratosthenesSequential` produces sorted results.
			primes := sieveOfEratosthenesSequential(tc.maxNum)
			if !reflect.DeepEqual(primes, tc.expected) {
				t.Errorf("sieveOfEratosthenesSequential(%d) = %v; want %v", tc.maxNum, primes, tc.expected)
			}
		})
	}
}

func TestFindPrimesWithSieve_Orchestration(t *testing.T) {
	testCases := []struct {
		name       string
		maxNum     int
		numWorkers int
		expected   []int // Expected results are the same, method differs
	}{
		{
			name:       "sieve primes up to 10 with 1 worker", // NumWorkers for sieve affects marking distribution
			maxNum:     10,
			numWorkers: 1,
			expected:   []int{2, 3, 5, 7},
		},
		{
			name:       "sieve primes up to 30 with 4 workers",
			maxNum:     30,
			numWorkers: 4,
			expected:   []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
		{
			name:       "sieve primes up to 50 with 2 workers",
			maxNum:     50,
			numWorkers: 2,
			expected:   []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47},
		},
		{
			name:       "sieve primes up to 1 (no primes) with 4 workers",
			maxNum:     1,
			numWorkers: 4,
			expected:   []int{},
		},
		// Add a slightly larger one if tests run fast enough
		{
			name:       "sieve primes up to 100 with 4 workers",
			maxNum:     100,
			numWorkers: 4,
			expected:   []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97},
		},
	}

	// Temporarily disable noisy logs for testing if needed
	// originalFlags := log.Flags()
	// log.SetOutput(io.Discard)
	// log.SetFlags(0)
	// defer func() {
	//  log.SetOutput(os.Stderr)
	//  log.SetFlags(originalFlags)
	// }()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// `expected` is already sorted. `findPrimesWithSieve` produces sorted results.
			primes := findPrimesWithSieve(tc.maxNum, tc.numWorkers)
			if !reflect.DeepEqual(primes, tc.expected) {
				t.Errorf("findPrimesWithSieve(%d, %d) = %v; want %v", tc.maxNum, tc.numWorkers, primes, tc.expected)
			}
		})
	}
}
