# Concurrent Prime Number Finder

A Go program that finds all prime numbers between 1 and a specified upper limit (e.g., 100,000) using Goroutines and Channels for parallelized computation.

## Features

- Efficient prime number checking (`isPrime` function).
- Basic concurrent pipeline using a worker goroutine, jobs channel, and results channel.
- Configurable number of workers and maximum number to check (via constants).
- Simple progress logging.

## Concurrency Model (4 Workers)

The system uses a fixed number of worker goroutines (e.g., 4) to parallelize prime checking.

1. The `findPrimes` function initializes:
    - A `jobs` channel (buffered): Numbers to be checked are sent here.
    - A `results` channel (buffered): Prime numbers found by workers are sent here.
    - A `sync.WaitGroup` to synchronize the completion of all worker goroutines.
2. **Worker Creation**: `numWorkers` (e.g., 4) `worker` goroutines are launched. Each worker is given access to the `jobs` and `results` channels, and the `WaitGroup`. `wg.Add(1)` is called for each worker.
3. **Job Dispatching**: A separate goroutine iterates from 1 to `maxNum`, sending each number to the `jobs` channel. After all numbers are sent, this goroutine closes the `jobs` channel. This signals to workers that no more jobs will be forthcoming.
4. **Worker Processing**: Each `worker` goroutine runs a loop:
    - It reads a number from the `jobs` channel (this blocks if the channel is empty).
    - It calls `isPrime` to check the number.
    - If `isPrime` returns `true`, the number (which is prime) is sent to the `results` channel.
    - When the `jobs` channel is closed by the dispatcher and the worker has processed all numbers it received from the channel, its loop terminates. The worker then calls `wg.Done()`.
5. **Result Collection and Synchronization**:
    - Another separate goroutine is launched. Its sole purpose is to wait for all worker goroutines to finish using `wg.Wait()`.
    - Once `wg.Wait()` unblocks (meaning all workers have called `wg.Done()`), this goroutine closes the `results` channel. This closure is the signal to the main result collector that no more primes will be sent.
6. **Final Aggregation**: The `findPrimes` function's main path of execution ranges over the `results` channel, appending each received prime number to a slice. This loop continues until the `results` channel is closed (by the goroutine in step 5).
7. The collected slice of primes is then sorted and returned.

This setup ensures that:

- Job distribution is handled by Go's channel mechanics (workers pull jobs as they become available).
- All numbers are processed.
- All primes are collected.
- The program terminates cleanly after all work is done.

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
