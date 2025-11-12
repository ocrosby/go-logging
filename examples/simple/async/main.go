// Async Logging Examples
// This example demonstrates asynchronous logging for high-performance applications

package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	fmt.Println("=== Async Logging Examples ===")

	// Basic async logging
	demonstrateBasicAsyncLogging()

	// High-throughput logging
	demonstrateHighThroughputLogging()

	// Async with graceful shutdown
	demonstrateGracefulShutdown()
}

func demonstrateBasicAsyncLogging() {
	fmt.Println("\n--- Basic Async Logging ---")

	// Create async output for high-performance logging
	underlying := logging.NewWriterOutput(os.Stdout)
	asyncOutput := logging.NewAsyncOutput(underlying, 100) // Queue size of 100
	defer asyncOutput.Close()

	// Create logger with async output
	logger := logging.NewEasyBuilder().
		JSON().
		Field("example", "basic_async").
		Build()

	// Log messages that will be processed asynchronously
	ctx := logging.NewContextWithTrace()

	for i := 0; i < 10; i++ {
		logger.InfoContext(ctx, "Async log message",
			"iteration", i,
			"timestamp", time.Now().UnixNano(),
			"worker_id", fmt.Sprintf("worker_%d", i%3),
		)

		// Small delay to show async processing
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("All messages submitted to async queue")
	time.Sleep(100 * time.Millisecond) // Allow async processing to complete
}

func demonstrateHighThroughputLogging() {
	fmt.Println("\n--- High Throughput Logging ---")

	// Create async output with larger queue
	underlying := logging.NewWriterOutput(os.Stdout)
	asyncOutput := logging.NewAsyncOutput(underlying, 1000)
	defer asyncOutput.Close()

	logger := logging.NewEasyBuilder().
		JSON().
		Field("example", "high_throughput").
		Build()

	// Simulate high-throughput logging from multiple goroutines
	var wg sync.WaitGroup
	numWorkers := 5
	messagesPerWorker := 20

	start := time.Now()

	for workerID := 0; workerID < numWorkers; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx := logging.NewContextWithTrace()
			ctx = logging.WithCorrelationID(ctx, fmt.Sprintf("worker_%d", id))

			for i := 0; i < messagesPerWorker; i++ {
				logger.InfoContext(ctx, "High-throughput message",
					"worker_id", id,
					"message_num", i,
					"payload", generatePayload(id, i),
					"timestamp", time.Now().UnixNano(),
				)

				// Simulate some work
				time.Sleep(1 * time.Millisecond)
			}
		}(workerID)
	}

	wg.Wait()
	elapsed := time.Since(start)

	logger.Info("High-throughput logging completed",
		"total_messages", numWorkers*messagesPerWorker,
		"workers", numWorkers,
		"duration_ms", elapsed.Milliseconds(),
		"messages_per_second", int(float64(numWorkers*messagesPerWorker)/elapsed.Seconds()),
	)

	// Allow async processing to complete
	time.Sleep(200 * time.Millisecond)
}

func demonstrateGracefulShutdown() {
	fmt.Println("\n--- Graceful Shutdown ---")

	// Create async output
	underlying := logging.NewWriterOutput(os.Stdout)
	asyncOutput := logging.NewAsyncOutput(underlying, 50)

	logger := logging.NewEasyBuilder().
		JSON().
		Field("example", "graceful_shutdown").
		Build()

	ctx := logging.NewContextWithTrace()

	// Start background logging
	done := make(chan bool)
	go func() {
		for i := 0; i < 25; i++ {
			logger.InfoContext(ctx, "Background logging",
				"iteration", i,
				"status", "processing",
			)
			time.Sleep(20 * time.Millisecond)
		}
		done <- true
	}()

	// Wait a bit then initiate shutdown
	time.Sleep(100 * time.Millisecond)

	logger.Info("Initiating graceful shutdown",
		"pending_messages", "will be processed",
	)

	// Stop async output - this waits for queued messages to be processed
	if err := asyncOutput.Stop(); err != nil {
		logger.Error("Error during shutdown",
			"error", err.Error(),
		)
	}

	// Wait for background work to complete
	<-done

	logger.Info("Graceful shutdown completed")
}

func generatePayload(workerID, messageNum int) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    fmt.Sprintf("user_%d", (workerID*100)+messageNum),
		"session_id": fmt.Sprintf("sess_%d_%d", workerID, messageNum),
		"action":     "data_processing",
		"details": map[string]interface{}{
			"processed_items": messageNum * 10,
			"cache_hits":      messageNum * 7,
			"cache_misses":    messageNum * 3,
		},
		"metadata": map[string]interface{}{
			"source":   "worker_pool",
			"priority": "normal",
			"retry":    messageNum%3 == 0,
		},
	}
}
