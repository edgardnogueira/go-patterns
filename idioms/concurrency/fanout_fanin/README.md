# Fan-Out/Fan-In Pattern

The Fan-Out/Fan-In pattern is a concurrency pattern that distributes work across multiple goroutines (fan-out) and then collects and merges the results into a single channel (fan-in). This pattern is ideal for parallel processing tasks where work can be divided into independent units.

## Pattern Overview

Fan-Out/Fan-In is particularly useful when:
- You have a large number of independent tasks
- Tasks can be processed in parallel
- You need to collect all results in one place

### Key Components

1. **Input Source**: Provides the work items to be processed
2. **Fan-Out Stage**: Distributes work to multiple goroutines
3. **Worker Functions**: Process individual work items
4. **Fan-In Stage**: Collects and combines results from all workers
5. **Output Channel**: Delivers combined results in a single stream

### Benefits

- **Increased Throughput**: Processes multiple items concurrently
- **Efficient Resource Utilization**: Makes optimal use of available CPU cores
- **Scalable Processing**: Easily adjust worker count based on workload
- **Simplified Collection**: Consolidates results into a single channel
- **Bounded Concurrency**: Controls how many goroutines are created

## Implementation Details

Our implementation provides:

- Generic functions working with any input/output types using Go generics
- Context awareness for proper cancellation
- Error handling and propagation
- Multiple convenience functions for different usage patterns
- Thread-safe result collection

## Common Use Cases

- **Data Processing**: Process large datasets in parallel
- **API Requests**: Make multiple concurrent API calls
- **File Operations**: Process multiple files simultaneously
- **Image Processing**: Apply transformations to multiple images
- **Search Operations**: Query multiple sources and combine results
- **Batch Processing**: Process items in batches with parallelism

## Example

See the [example directory](./example) for a complete working example demonstrating the Fan-Out/Fan-In pattern in action, showcasing a web scraper that processes multiple URLs concurrently.

## Best Practices

- Choose an appropriate worker count based on your workload (CPU vs I/O bound)
- Consider using buffered channels for smoother flow of data
- Always handle errors from worker goroutines
- Use context cancellation to gracefully stop all operations
- Be mindful of memory usage when processing large datasets
- Consider rate limiting if interacting with external resources
