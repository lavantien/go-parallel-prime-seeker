package main

import (
	"log"
	"math"
	"sort" // For final sorting
	"sync"
	"time"
)

const (
	MaxNumberForSieveGlobal           = 200_000_000
	NumSieveWorkersGlobal             = 4
	SieveProgressReportIntervalGlobal = 1_000_000 // Approximate interval for logging prime collection
	// SegmentSizeInNumbers defines how many numbers each segment will handle.
	// A larger size reduces overhead but increases memory per worker and cache pressure.
	// 524288 numbers = 64KB bitset per worker (524288 / 8 bytes)
	SegmentSizeInNumbers = 524288
)

// --- Bitset Helper Functions (for segment bitsets) ---

// markBitSegment sets the bit corresponding to indexInSegment in the segmentBitset.
// indexInSegment is (number - segmentLow).
func markBitSegment(indexInSegment int, segmentBitset []byte) {
	byteIndex := indexInSegment / 8
	bitOffset := uint(indexInSegment % 8)
	// No bounds check here, assuming segmentBitset is correctly sized for indexInSegment
	segmentBitset[byteIndex] |= (1 << bitOffset)
}

// isBitMarkedSegment checks if the bit corresponding to indexInSegment is set.
// indexInSegment is (number - segmentLow).
func isBitMarkedSegment(indexInSegment int, segmentBitset []byte) bool {
	byteIndex := indexInSegment / 8
	bitOffset := uint(indexInSegment % 8)
	// No bounds check here
	return (segmentBitset[byteIndex] & (1 << bitOffset)) != 0
}

// --- Sequential Sieve for Base Primes (Optimized) ---
func sieveOfEratosthenesSequentialBase(maxNumInternal int) []int {
	if maxNumInternal < 2 {
		return []int{}
	}
	// isComposite[i] is true if odd number i is composite.
	// Max index needed is maxNumInternal.
	isComposite := make([]bool, maxNumInternal+1)

	// Sieving for odd primes:
	// Start p from 3. Mark multiples p*p, p*(p+2), p*(p+4), ...
	// which are p*p, p*p + 2p, p*p + 4p, ...
	for p := 3; int64(p)*int64(p) <= int64(maxNumInternal); p += 2 { // Loop up to sqrt for marking
		if !isComposite[p] {
			for multiple := p * p; multiple <= maxNumInternal; multiple += (2 * p) { // Only mark odd multiples
				isComposite[multiple] = true
			}
		}
	}

	primesList := make([]int, 0)
	if maxNumInternal >= 2 {
		primesList = append(primesList, 2) // Add 2, the only even prime
	}
	for i := 3; i <= maxNumInternal; i += 2 { // Collect by iterating up to maxNumInternal
		if !isComposite[i] {
			primesList = append(primesList, i)
		}
	}
	return primesList
}

// SegmentTask defines a piece of work for a worker.
type SegmentTask struct {
	low, high int // Range [low, high] to sieve
	// id        int // Optional: for ordered collection if strictly needed before final sort
}

// SegmentResult holds primes found in a segment.
type SegmentResult struct {
	primes []int
	// id     int
}

// segmentedSieveWorker processes segments sent via the tasks channel.
func segmentedSieveWorker(
	workerID int,
	tasks <-chan SegmentTask,
	results chan<- SegmentResult,
	basePrimes []int, // Primes up to sqrt(maxNum)
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	// log.Printf("Segmented Worker %d: Starting", workerID)

	for task := range tasks {
		segmentNumCount := task.high - task.low + 1
		// log.Printf("Segmented Worker %d: Task for range [%d, %d] (%d numbers)", workerID, task.low, task.high, segmentNumCount)

		// Create a bitset for the current segment [task.low, task.high].
		// The bitset index `k` corresponds to the number `task.low + k`.
		segmentBitsetLen := (segmentNumCount + 7) / 8
		segmentBitset := make([]byte, segmentBitsetLen) // Initialized to all zeros (all potentially prime)

		for _, p := range basePrimes {
			// Optimization: if p*p > task.high, then this p and subsequent larger primes
			// cannot have any multiples p*p (or larger) in this segment if we consider
			// that smaller multiples (p*k where k<p) would have been marked by smaller primes.
			// This prime p starts marking at p*p. If p*p is beyond the segment, it won't mark anything in this segment.
			if int64(p)*int64(p) > int64(task.high) {
				break
			}

			// Calculate the first multiple of p that is >= task.low
			// startMultipleInP = ceil(task.low / p) * p
			startMultipleInP := ((task.low + p - 1) / p) * p

			// The actual marking should start from max(p*p, startMultipleInP)
			// because multiples p*k where k < p have already been handled by prime k.
			actualStartMarking := startMultipleInP
			if actualStartMarking < p*p {
				actualStartMarking = p * p
			}

			// Mark multiples of p within the segment [task.low, task.high]
			for multiple := actualStartMarking; multiple <= task.high; multiple += p {
				// The index in segmentBitset is (multiple - task.low)
				if multiple >= task.low { // Ensure multiple is within current segment (should be due to loop condition)
					markBitSegment(multiple-task.low, segmentBitset)
				}
			}
		}

		// Collect primes from this segment's bitset
		segmentPrimes := make([]int, 0, segmentNumCount/10) // Rough pre-allocation
		for i := 0; i < segmentNumCount; i++ {
			currentNum := task.low + i
			if currentNum < 2 { // Primes are >= 2
				continue
			}
			if !isBitMarkedSegment(i, segmentBitset) {
				segmentPrimes = append(segmentPrimes, currentNum)
			}
		}
		results <- SegmentResult{primes: segmentPrimes}
	}
	// log.Printf("Segmented Worker %d: Exiting", workerID)
}

// findPrimesWithSegmentedSieve implements a parallel segmented Sieve of Eratosthenes.
func findPrimesWithSegmentedSieve(maxNum int, numWorkers int) []int {
	if maxNum < 2 {
		return []int{}
	}
	overallStartTime := time.Now()

	// --- Phase 1: Find base primes up to sqrt(maxNum) ---
	sqrtMaxNum := int(math.Sqrt(float64(maxNum)))
	basePrimesStartTime := time.Now()
	log.Printf("Segmented Sieve: Finding base primes up to %d", sqrtMaxNum)
	basePrimes := sieveOfEratosthenesSequentialBase(sqrtMaxNum)
	log.Printf("Segmented Sieve: Found %d base primes in %s", len(basePrimes), time.Since(basePrimesStartTime))

	if len(basePrimes) == 0 && maxNum >= 2 && sqrtMaxNum >= 2 {
		log.Printf("Segmented Sieve: Warning - No base primes found though sqrtMaxNum (%d) >= 2. This implies maxNum is small (e.g., <4). Segments will handle prime finding.", sqrtMaxNum)
	}

	tasks := make(chan SegmentTask, numWorkers)     // Channel for workers to pick up segments
	results := make(chan SegmentResult, numWorkers) // Channel for workers to send back found primes

	// --- Start Worker Goroutines ---
	var wgWorkers sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wgWorkers.Add(1)
		go segmentedSieveWorker(w, tasks, results, basePrimes, &wgWorkers)
	}

	// --- Phase 2: Dispatch Segment Tasks ---
	dispatchTime := time.Now()
	// Calculate number of segments first for the collector's WaitGroup
	numDispatchedSegments := 0
	low := 0
	for low <= maxNum {
		numDispatchedSegments++
		low += SegmentSizeInNumbers
	}

	go func() {
		defer close(tasks) // Close tasks channel when all tasks are sent
		currentLow := 0
		for currentLow <= maxNum {
			currentHigh := currentLow + SegmentSizeInNumbers - 1
			if currentHigh > maxNum {
				currentHigh = maxNum
			}
			if currentLow > currentHigh { // Should not happen if loop condition is currentLow <= maxNum
				break
			}
			tasks <- SegmentTask{low: currentLow, high: currentHigh}
			currentLow += SegmentSizeInNumbers
		}
		log.Printf("Segmented Sieve: All %d segment tasks dispatched in %s", numDispatchedSegments, time.Since(dispatchTime))
	}()

	// --- Phase 3: Collect and Combine Results ---
	collectionStartTime := time.Now()
	var wgCollector sync.WaitGroup
	wgCollector.Add(numDispatchedSegments) // Expect one result per dispatched segment

	intermediateCollectedPrimes := make([][]int, 0, numDispatchedSegments)
	collectedPrimesCount := 0
	uniquePrimesFound := 0 // This will be len(finalPrimes) after sort & unique

	// Goroutine to collect results and manage wgCollector
	go func() {
		for result := range results {
			intermediateCollectedPrimes = append(intermediateCollectedPrimes, result.primes)
			collectedPrimesCount += len(result.primes) // Sum of primes in all segments (before sort/unique)
			// Simple progress, not tied to SieveProgressReportIntervalGlobal yet
			// log.Printf("Segmented Sieve: Collected segment result (approx %d primes so far)", collectedPrimesCount)
			wgCollector.Done()
		}
	}()

	wgCollector.Wait() // Wait for all segment results to be collected
	close(results)     // Close results channel as all results are in

	log.Printf("Segmented Sieve: All segment results (%d segments) collected in %s. Raw primes collected: %d",
		numDispatchedSegments, time.Since(collectionStartTime), collectedPrimesCount)

	// --- Final Assembly and Sorting ---
	assemblyStartTime := time.Now()
	// Estimate capacity for the final list of primes
	finalCapacity := 0
	if maxNum > 1 {
		logMax := math.Log(float64(maxNum))
		if logMax > 0 {
			finalCapacity = int(float64(maxNum) / logMax) // Prime Number Theorem approximation
		}
	}
	if finalCapacity <= 0 { // Fallback for small maxNum or if PNT estimate is off
		finalCapacity = collectedPrimesCount / 2 // A rough heuristic
		if finalCapacity < 10 {
			finalCapacity = 10
		}
	}

	finalPrimes := make([]int, 0, finalCapacity)
	for _, segmentPrimes := range intermediateCollectedPrimes {
		finalPrimes = append(finalPrimes, segmentPrimes...)
	}

	// Sort the combined list to ensure primes are in order.
	// This step also implicitly handles uniqueness if primes were somehow redundantly generated
	// by different segments (which they shouldn't be with disjoint segments).
	sort.Ints(finalPrimes)
	uniquePrimesFound = len(finalPrimes) // After sort, len gives unique prime count if no duplicates

	log.Printf("Segmented Sieve: Primes combined and sorted in %s. Total unique primes found: %d",
		time.Since(assemblyStartTime), uniquePrimesFound)
	log.Printf("Segmented Sieve: Total time for findPrimesWithSegmentedSieve: %s", time.Since(overallStartTime))

	// Ensure workers are fully done (though they should be if tasks and results are closed)
	wgWorkers.Wait()

	// Progress reporting is a bit coarse here. Could be integrated into collection.
	// For instance, print every X primes appended to finalPrimes during the append loop,
	// but that might slow it down. The current logging provides phase timings.

	return finalPrimes
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("Concurrent Prime Finder (Segmented Sieve Version) - Starting")

	maxNumber := MaxNumberForSieveGlobal
	numWorkers := NumSieveWorkersGlobal

	if numWorkers <= 0 {
		log.Println("Number of workers must be positive. Defaulting to 1.")
		numWorkers = 1
	}
	if SegmentSizeInNumbers <= 0 {
		log.Fatal("SegmentSizeInNumbers must be positive.")
	}

	log.Printf("Finding primes up to %d using %d workers (Segmented Sieve method).\nSegment size: %d numbers.\n",
		maxNumber, numWorkers, SegmentSizeInNumbers)

	primes := findPrimesWithSegmentedSieve(maxNumber, numWorkers)

	log.Printf("Found %d prime numbers up to %d.\n", len(primes), maxNumber)

	// Example verification (optional)
	// if maxNumber <= 200 {
	// 	log.Printf("Primes (all): %v\n", primes)
	// } else if len(primes) > 20 {
	// 	log.Printf("First 10 primes: %v\n", primes[:10])
	// 	log.Printf("Last 10 primes: %v\n", primes[len(primes)-10:])
	// }

	// Simple check based on known counts for small N
	// if maxNumber == 200_000_000 {
	//    expected_pi_200M := 11078937 // From online sources pi(2*10^8)
	//    if len(primes) == expected_pi_200M {
	//        log.Printf("Verification: Prime count matches expected for N=200,000,000 (%d).", expected_pi_200M)
	//    } else {
	//        log.Printf("Verification: Prime count MISMATCH for N=200,000,000. Expected: %d, Got: %d", expected_pi_200M, len(primes))
	//    }
	// }

	log.Println("Concurrent Prime Finder (Segmented Sieve Version) - Finished")
}
