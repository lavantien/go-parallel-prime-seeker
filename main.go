package main

import (
	"io"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	sqrtN := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrtN; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	// log.Printf("Worker %d starting\n", id)
	for num := range jobs {
		if isPrime(num) {
			results <- num
		}
	}
	// log.Printf("Worker %d finished\n", id)
}

// Constants
const (
	DefaultMaxNumber     = 10000000
	DefaultNumWorkers    = 8
	ProgressReportPrimes = 10000 // Report every 1000 primes found
)

// findPrimes orchestrates the prime finding process with multiple workers.
func findPrimes(maxNum int, numWorkers int) []int {
	jobs := make(chan int, maxNum)
	results := make(chan int, maxNum)
	var wg sync.WaitGroup

	// Launch numWorkers goroutines.
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs (numbers to check) in a separate goroutine.
	go func() {
		for i := 1; i <= maxNum; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Collector goroutine: Waits for all workers to finish, then closes results.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from the results channel until it's closed.
	primes := make([]int, 0) // Initialize as empty, non-nil slice

	foundCount := 0
	startTime := time.Now()

	// log.Printf("Collecting results...\n") // Keep or remove based on desired test log verbosity

	for prime := range results {
		primes = append(primes, prime)
		foundCount++
		// Progress reporting logic can be here or simplified for tests
		// For tests, often less logging is better unless debugging.
		// The current progress logging is fine for main execution.
		// It might be slightly noisy for `go test -v` if many test cases run findPrimes.
		// For now, let's assume it's acceptable.
		if foundCount%ProgressReportPrimes == 0 && log.Default().Writer() != io.Discard { // Check if logging is enabled
			log.Printf("Found %d primes so far... (last prime: %d)\n", foundCount, prime)
		}
	}

	if log.Default().Writer() != io.Discard {
		elapsedTime := time.Since(startTime)
		log.Printf("Collected all results in %s.\n", elapsedTime)
	}

	sort.Ints(primes) // Sort for consistent output and easier testing.
	return primes
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds) // Add timestamps to log
	log.Println("Concurrent Prime Finder - Starting")

	maxNumber := DefaultMaxNumber
	numWorkers := DefaultNumWorkers

	log.Printf("Finding primes up to %d using %d workers.\n", maxNumber, numWorkers)

	primes := findPrimes(maxNumber, numWorkers)

	log.Printf("Found %d prime numbers up to %d.\n", len(primes), maxNumber)
	log.Println("Concurrent Prime Finder - Finished")
}
