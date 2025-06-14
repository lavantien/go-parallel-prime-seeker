package main

import (
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

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Job dispatcher goroutine
	go func() {
		// For progress based on numbers dispatched:
		// dispatchedCount := 0
		// lastReportedDispatched := 0
		for i := 1; i <= maxNum; i++ {
			jobs <- i
			// dispatchedCount++
			// if dispatchedCount-lastReportedDispatched >= ProgressReportNumbers {
			//  log.Printf("Dispatched %d numbers for checking...\n", dispatchedCount)
			//  lastReportedDispatched = dispatchedCount
			// }
		}
		close(jobs)
		// log.Printf("All %d numbers dispatched for checking.\n", dispatchedCount)
	}()

	// Collector goroutine for closing results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	var primes []int
	foundCount := 0
	startTime := time.Now()

	log.Printf("Collecting results...\n")

	for prime := range results {
		primes = append(primes, prime)
		foundCount++
		if foundCount%ProgressReportPrimes == 0 {
			// This will print progress for every 'ProgressReportPrimes' found
			log.Printf("Found %d primes so far... (last prime: %d)\n", foundCount, prime)
		}
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Collected all results in %s.\n", elapsedTime)

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
