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

// sieveOfEratosthenesSequential finds all primes up to maxNum using a sequential sieve.
func sieveOfEratosthenesSequential(maxNum int) []int {
	if maxNum < 2 {
		return []int{}
	}

	// Create a boolean slice `isNotPrime` indexed from 0 to maxNum.
	// isNotPrime[i] is true if i is NOT prime.
	// Initialize all to false (meaning potentially prime).
	isNotPrime := make([]bool, maxNum+1)
	isNotPrime[0] = true // 0 is not prime
	isNotPrime[1] = true // 1 is not prime

	// Iterate from 2 up to sqrt(maxNum)
	sqrtMax := int(math.Sqrt(float64(maxNum)))
	for p := 2; p <= sqrtMax; p++ {
		// If p is prime (i.e., isNotPrime[p] is false)
		if !isNotPrime[p] {
			// Mark all multiples of p as not prime
			// Start marking from p*p, as smaller multiples
			// would have been marked by smaller primes.
			for multiple := p * p; multiple <= maxNum; multiple += p {
				isNotPrime[multiple] = true
			}
		}
	}

	// Collect primes
	primes := make([]int, 0)
	for i := 2; i <= maxNum; i++ {
		if !isNotPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

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
