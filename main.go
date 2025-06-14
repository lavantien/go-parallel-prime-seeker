package main

import (
	"fmt"
	"math"
)

// isPrime checks if a number n is prime.
// Optimizations:
// 1. Numbers less than 2 are not prime.
// 2. 2 is the only even prime number.
// 3. Check divisibility only up to sqrt(n).
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

func main() {
	fmt.Println("Concurrent Prime Finder")
	// Example usage (will be removed later)
	fmt.Printf("Is 7 prime? %v\n", isPrime(7))
	fmt.Printf("Is 10 prime? %v\n", isPrime(10))
}
