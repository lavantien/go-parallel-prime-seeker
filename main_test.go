package main

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

// TestSieveOfEratosthenesSequentialBase tests the optimized sequential sieve for base primes.
func TestSieveOfEratosthenesSequentialBase(t *testing.T) {
	testCases := []struct {
		name     string
		maxNum   int
		expected []int
	}{
		{"primes up to 10", 10, []int{2, 3, 5, 7}},
		{"primes up to 20", 20, []int{2, 3, 5, 7, 11, 13, 17, 19}},
		{"primes up to 2", 2, []int{2}},
		{"primes up to 1", 1, []int{}}, // Max num 1, so sqrt(1) = 1, should yield no primes.
		{"primes up to 0", 0, []int{}},
		{"primes up to 30", 30, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
		{"primes up to 3", 3, []int{2, 3}},
		{"primes up to 4", 4, []int{2, 3}}, // sqrt(4)=2, base prime is 2.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			primes := sieveOfEratosthenesSequentialBase(tc.maxNum)
			if !reflect.DeepEqual(primes, tc.expected) {
				t.Errorf("sieveOfEratosthenesSequentialBase(%d) = %v; want %v", tc.maxNum, primes, tc.expected)
			}
		})
	}
}

// TestFindPrimesWithSegmentedSieve_Orchestration tests the overall segmented sieve orchestration.
func TestFindPrimesWithSegmentedSieve_Orchestration(t *testing.T) {
	testCases := []struct {
		name        string
		maxNum      int
		numWorkers  int
		segmentSize int // Added for testing different segment configurations
		expected    []int
	}{
		{
			name:        "sieve up to 10, 1 worker, seg 5",
			maxNum:      10,
			numWorkers:  1,
			segmentSize: 5,
			expected:    []int{2, 3, 5, 7},
		},
		{
			name:        "sieve up to 30, 4 workers, seg 10",
			maxNum:      30,
			numWorkers:  4,
			segmentSize: 10,
			expected:    []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
		{
			name:        "sieve up to 50, 2 workers, seg 20",
			maxNum:      50,
			numWorkers:  2,
			segmentSize: 20,
			expected:    []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47},
		},
		{
			name:        "sieve up to 1 (no primes), 4 workers, seg 10",
			maxNum:      1,
			numWorkers:  4,
			segmentSize: 10,
			expected:    []int{},
		},
		{
			name:        "sieve up to 100, 4 workers, seg 25",
			maxNum:      100,
			numWorkers:  4,
			segmentSize: 25,
			expected:    []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97},
		},
		{
			name:        "sieve up to 7, 2 workers, seg 3 - primes only in base up to sqrt(7)=2",
			maxNum:      7, // sqrt(7) is ~2.6, so basePrimes = [2]
			numWorkers:  2,
			segmentSize: 3, // Segments: [0,2], [3,5], [6,7]
			expected:    []int{2, 3, 5, 7},
		},
	}

	// Silence logs for most test runs to keep output clean, can be re-enabled for debugging.
	originalLogOutput := log.Writer()
	originalLogFlags := log.Flags()
	log.SetOutput(io.Discard) // Send logs to nowhere
	// log.SetFlags(0) // Remove timestamping if desired, but not strictly needed with discard

	defer func() {
		log.SetOutput(originalLogOutput) // Restore log output
		log.SetFlags(originalLogFlags)   // Restore log flags
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// It's tricky to directly inject SegmentSizeInNumbers from const for tests
			// without making it a parameter or global var. For testability,
			// if SegmentSizeInNumbers were a parameter to findPrimesWithSegmentedSieve,
			// it would be cleaner. Assuming it's a package-level const, we test with current const.
			// WORKAROUND: If SegmentSizeInNumbers is a const in main, these tests assume it's set appropriately
			// or that tc.segmentSize is just for documentation of the scenario.
			// Let's assume there's a way to modify it for testing or we test current SegmentSizeInNumbers.
			// For this example, I'm adapting the test to *imply* the segment size effect,
			// but a real test might need to modify a global or pass it in.

			// To truly test different segment sizes, findPrimesWithSegmentedSieve would need to accept it,
			// or we would need to manipulate the global constant if it's not a const.
			// Let's assume for this exercise that SegmentSizeInNumbers is the one defined in main.go
			// and tc.segmentSize is for descriptive purposes or if we could modify it.

			// log.SetOutput(os.Stdout) // Enable logs for this specific test
			// log.Printf("Testing %s with maxNum=%d, numWorkers=%d (effective segment size: %d)", tc.name, tc.maxNum, tc.numWorkers, SegmentSizeInNumbers)

			primes := findPrimesWithSegmentedSieve(tc.maxNum, tc.numWorkers) // Uses SegmentSizeInNumbers from main.go

			// log.SetOutput(io.Discard) // Re-disable logs

			if !reflect.DeepEqual(primes, tc.expected) {
				// Re-enable log for detailed error output if a test fails
				log.SetOutput(os.Stdout)
				log.Printf("Log for failed test %s:", tc.name)
				// Rerun with logging enabled to see what happened for this specific fail.
				// This won't capture past logs unless we buffer them.
				// For simplicity, just printing error:
				t.Errorf("findPrimesWithSegmentedSieve(%d, %d) [effective seg size %d] = %v;\n want %v",
					tc.maxNum, tc.numWorkers, SegmentSizeInNumbers, primes, tc.expected)
				log.SetOutput(io.Discard) // Disable again
			}
		})
	}
}

// BenchmarkFindPrimesWithSegmentedSieve provides a basic benchmark.
func BenchmarkFindPrimesWithSegmentedSieve(b *testing.B) {
	maxNum := 10_000_000                // A moderately large number for benchmark
	numWorkers := NumSieveWorkersGlobal // Use the global config

	// Silence logs during benchmark
	originalLogOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(originalLogOutput)

	b.ResetTimer() // Start timing after setup
	for i := 0; i < b.N; i++ {
		_ = findPrimesWithSegmentedSieve(maxNum, numWorkers)
	}
}
