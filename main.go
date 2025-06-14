package main

import (
	"log"
	"math" // Still useful for base primes or if order matters before returning
	"sync"
	"time"
)

const (
	MaxNumberForSieve           = 200_000_000
	NumSieveWorkers             = 4         // Requirement
	SieveProgressReportInterval = 1_000_000 // Report progress by numbers sieved or primes found
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

// Sequential Sieve
func sieveOfEratosthenesSequential(maxNum int) []int {
	if maxNum < 2 {
		return []int{}
	}
	isNotPrime := make([]bool, maxNum+1)
	isNotPrime[0] = true
	isNotPrime[1] = true
	sqrtMax := int(math.Sqrt(float64(maxNum)))
	for p := 2; p <= sqrtMax; p++ {
		if !isNotPrime[p] {
			for multiple := p * p; multiple <= maxNum; multiple += p {
				isNotPrime[multiple] = true
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= maxNum; i++ {
		if !isNotPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// findPrimesWithSieve implements a parallel Sieve of Eratosthenes.
func findPrimesWithSieve(maxNum int, numWorkers int) []int {
	if maxNum < 2 {
		return []int{}
	}

	startTime := time.Now()

	// Phase 1: Find base primes up to sqrt(maxNum) using a sequential sieve.
	// This is efficient because sqrt(maxNum) is relatively small.
	sqrtMaxNum := int(math.Sqrt(float64(maxNum)))
	log.Printf("Sieve: Finding base primes up to %d\n", sqrtMaxNum)
	basePrimes := sieveOfEratosthenesSequential(sqrtMaxNum)
	log.Printf("Sieve: Found %d base primes in %s\n", len(basePrimes), time.Since(startTime))

	// Phase 2: Initialize the main sieve array (bitset for memory efficiency with large N, but bool slice for simplicity here).
	// `isNotPrime[i] = true` means 'i' is composite or 0 or 1.
	sieveTime := time.Now()
	isNotPrime := make([]bool, maxNum+1)
	isNotPrime[0] = true
	isNotPrime[1] = true

	// Phase 3: Parallel marking of multiples.
	var wg sync.WaitGroup

	// Distribute basePrimes among workers for marking.
	// Each worker handles a contiguous block of basePrimes.
	primesPerWorker := (len(basePrimes) + numWorkers - 1) / numWorkers

	log.Printf("Sieve: Starting parallel marking with %d workers. Total base primes: %d, per worker approx: %d\n",
		numWorkers, len(basePrimes), primesPerWorker)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Determine the range of basePrimes for this worker
			startIdx := workerID * primesPerWorker
			endIdx := (workerID + 1) * primesPerWorker
			if endIdx > len(basePrimes) {
				endIdx = len(basePrimes)
			}

			// log.Printf("Sieve Worker %d: processing base primes from index %d to %d\n", workerID, startIdx, endIdx-1)

			for i := startIdx; i < endIdx; i++ {
				p := basePrimes[i]
				// Mark multiples of p. Start marking from p*p.
				// All numbers p*k where k < p would have already been marked
				// by a prime factor smaller than p.
				for multiple := int64(p) * int64(p); multiple <= int64(maxNum); multiple += int64(p) {
					if multiple < 0 {
						break
					} // Overflow check though unlikely with bool slice index
					isNotPrime[multiple] = true
				}
			}
			// log.Printf("Sieve Worker %d: finished processing its base primes.\n", workerID)
		}(w)
	}

	wg.Wait() // Wait for all workers to finish marking
	log.Printf("Sieve: Parallel marking completed in %s\n", time.Since(sieveTime))

	// Phase 4: Collect results
	collectTime := time.Now()
	primes := make([]int, 0, maxNum/10) // Pre-allocate with a rough estimate

	// Progress reporting setup
	foundCount := 0

	for i := 2; i <= maxNum; i++ {
		if !isNotPrime[i] {
			primes = append(primes, i)
			foundCount++
			if foundCount%SieveProgressReportInterval == 0 {
				log.Printf("Sieve: Collected %d primes so far (last: %d)...\n", foundCount, i)
			}
		}
	}
	log.Printf("Sieve: Result collection completed in %s. Total primes: %d\n", time.Since(collectTime), len(primes))
	log.Printf("Sieve: Total time for findPrimesWithSieve: %s\n", time.Since(startTime))

	return primes
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("Concurrent Prime Finder (Sieve Version) - Starting")

	maxNumber := MaxNumberForSieve
	numWorkers := NumSieveWorkers

	log.Printf("Finding primes up to %d using %d workers (Sieve method).\n", maxNumber, numWorkers)

	primes := findPrimesWithSieve(maxNumber, numWorkers)

	log.Printf("Found %d prime numbers up to %d.\n", len(primes), maxNumber)
	// Optionally print some primes if range is small
	// if maxNumber <= 100 {
	//  log.Printf("Primes: %v\n", primes)
	// }
	log.Println("Concurrent Prime Finder (Sieve Version) - Finished")
}
