# Concurrent Prime Number Finder

A Go program that finds all prime numbers between 1 and a specified upper limit (e.g., 100,000) using Goroutines and Channels for parallelized computation.

## Features

- Efficient prime number checking (`isPrime` function).
- Efficient prime number checking (`isPrime` function).
- Basic concurrent pipeline using a worker goroutine, jobs channel, and results channel.

## Concurrency Model (Initial)

1. The `main` function (or a dedicated orchestrator like `findPrimes`) creates:
    - A `jobs` channel: To send numbers that need to be checked for primality.
    - A `results` channel: To receive prime numbers found by workers.
    - A `sync.WaitGroup` to track completion of worker goroutines.
2. A (currently single) `worker` goroutine is launched.
3. Numbers from 1 to `maxNum` are sent to the `jobs` channel by a dedicated goroutine. Once all numbers are sent, the `jobs` channel is closed.
4. The `worker` goroutine reads numbers from the `jobs` channel. For each number:
    - It calls `isPrime`.
    - If the number is prime, it's sent to the `results` channel.
5. When the `jobs` channel is closed and emptied, the `worker` goroutine finishes and calls `wg.Done()`.
6. A separate goroutine waits for all workers to complete (using `wg.Wait()`) and then closes the `results` channel.
7. The `main` function collects prime numbers from the `results` channel until it's closed.

## Requirements

- Go (version 1.x)

## How to Run

```bash
go run main.go
```

### Running Tests

```bash
go test -v
```
