package main

import (
	"log" // Changed to log for potential progress later
	"math"
	"sort"
	"sync"
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
	// log.Printf("Worker %d starting\n", id) // Can be noisy, make optional
	for num := range jobs {
		if isPrime(num) {
			results <- num
		}
	}
	// log.Printf("Worker %d finished\n", id) // Can be noisy
}

// Constants
const (
	DefaultMaxNumber  = 100000
	DefaultNumWorkers = 4
)

// findPrimes orchestrates the prime finding process with multiple workers.
func findPrimes(maxNum int, numWorkers int) []int {
	jobs := make(chan int, maxNum)    // Buffered channel for jobs
	results := make(chan int, maxNum) // Buffered channel for results
	var wg sync.WaitGroup

	// Launch numWorkers goroutines.
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs (numbers to check) in a separate goroutine.
	// This allows the main flow to proceed to collecting results.
	go func() {
		for i := 1; i <= maxNum; i++ {
			jobs <- i
		}
		close(jobs) // Close jobs channel when all numbers are sent.
	}()

	// Collector goroutine: Waits for all workers to finish, then closes results.
	// This is crucial. If wg.Wait() and close(results) were in the main `findPrimes`
	// goroutine before the results collection loop, it would deadlock because
	// it would wait for workers that might be blocked trying to send to `results`.
	go func() {
		wg.Wait()
		close(results) // Close results channel after all workers are done.
	}()

	// Collect results from the results channel until it's closed.
	var primes []int
	for prime := range results {
		primes = append(primes, prime)
	}

	sort.Ints(primes) // Sort for consistent output and easier testing.
	return primes
}

func main() {
	log.Println("Concurrent Prime Finder - Starting")

	maxNumber := DefaultMaxNumber
	numWorkers := DefaultNumWorkers // Use 4 workers as per requirement

	log.Printf("Finding primes up to %d using %d workers.\n", maxNumber, numWorkers)

	primes := findPrimes(maxNumber, numWorkers)

	log.Printf("Found %d prime numbers up to %d.\n", len(primes), maxNumber)
	// Optionally print all primes, but it's a lot for 100,000
	// if maxNumber <= 100 {
	//  log.Printf("Primes: %v\n", primes)
	// }
	log.Println("Concurrent Prime Finder - Finished")
}
