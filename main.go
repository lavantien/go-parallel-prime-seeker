package main

import (
	"fmt"
	"math"
	"sort" // For consistent output in tests/main
	"sync"
)

// isPrime function (from previous step)
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

// worker function
// It receives numbers from the jobs channel.
// If a number is prime, it sends it to the results channel.
// wg.Done() is called when the jobs channel is closed and all jobs are processed.
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this worker is done when the function returns
	fmt.Printf("Worker %d starting\n", id)
	for num := range jobs {
		// fmt.Printf("Worker %d processing %d\n", id, num) // Optional: for detailed logging
		if isPrime(num) {
			results <- num
		}
	}
	fmt.Printf("Worker %d finished\n", id)
}

// findPrimes orchestrates the prime finding process.
// For now, it uses a single worker.
func findPrimes(maxNum int, numWorkers int) []int {
	jobs := make(chan int, maxNum)    // Buffer size can be tuned
	results := make(chan int, maxNum) // Buffer size can be tuned
	var wg sync.WaitGroup

	// Launch worker(s)
	// For this commit, we'll explicitly launch just one for simplicity in understanding
	// but design the function to take numWorkers.
	// In the next step, we'll loop to launch `numWorkers`.
	wg.Add(1)
	go worker(1, jobs, results, &wg) // For now, only one worker with ID 1

	// Send jobs (numbers to check)
	// This should be in a separate goroutine so it doesn't block closing `results` later.
	go func() {
		for i := 1; i <= maxNum; i++ {
			jobs <- i
		}
		close(jobs) // Signal workers that no more jobs are coming
	}()

	// Wait for all workers to finish, then close the results channel.
	// This must be in a separate goroutine to avoid deadlock.
	// The main goroutine (findPrimes) needs to collect results,
	// so it can't block on wg.Wait() before starting to collect.
	go func() {
		wg.Wait()
		close(results) // Signal that no more results will be sent
	}()

	// Collect results
	var primes []int
	for prime := range results {
		primes = append(primes, prime)
	}

	sort.Ints(primes) // Sort for consistent output and easier testing
	return primes
}

func main() {
	fmt.Println("Concurrent Prime Finder")
	maxNumber := 30 // Small range for now
	numWorkers := 1 // Single worker for this step

	primes := findPrimes(maxNumber, numWorkers)

	fmt.Printf("Prime numbers up to %d: %v\n", maxNumber, primes)
	fmt.Printf("Found %d primes.\n", len(primes))
}
