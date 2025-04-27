package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/edgardnogueira/go-patterns/idioms/concurrency/fanout_fanin"
)

// WeatherData represents weather information for a location
type WeatherData struct {
	Location    string
	Temperature float64
	WindSpeed   float64
	Humidity    int
	Timestamp   time.Time
}

// WeatherResult contains processed weather data
type WeatherResult struct {
	Location      string
	Analysis      string
	WarningLevel  int
	ProcessedBy   int
	ProcessingTime time.Duration
}

// List of locations to fetch weather data for
var locations = []string{
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
	"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
	"Austin", "Jacksonville", "Fort Worth", "Columbus", "San Francisco",
	"Charlotte", "Indianapolis", "Seattle", "Denver", "Washington",
	"Boston", "El Paso", "Nashville", "Detroit", "Portland",
}

func main() {
	fmt.Println("=== Fan-Out/Fan-In Pattern - Weather Data Processing System ===")
	fmt.Println("---------------------------------------------------------------")
	
	// Parse worker count from args or use default
	workerCount := 4
	if len(os.Args) > 1 {
		if count, err := strconv.Atoi(os.Args[1]); err == nil && count > 0 {
			workerCount = count
		}
	}
	
	// Create a context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutdown signal received. Canceling operations...")
		cancel()
	}()
	
	fmt.Printf("Starting weather data processing with %d concurrent workers\n", workerCount)
	fmt.Println("Press Ctrl+C to stop the program")
	fmt.Println("---------------------------------------------------------------")
	
	// Run the main processing loop
	processWeatherData(ctx, workerCount)
}

// processWeatherData runs the main processing loop
func processWeatherData(ctx context.Context, workerCount int) {
	// Create a ticker for fetching data periodically
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	// Statistics counters
	var (
		mu            sync.Mutex
		processedCount int
		warningCount   int
		totalTime      time.Duration
	)
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("\nFetching weather data for", len(locations), "locations...")
			
			// 1. Fetch weather data (simulated)
			weatherDataChan := fetchWeatherData(ctx, locations)
			
			// 2. Process weather data using Fan-Out/Fan-In pattern
			resultChan, errChan := fanout_fanin.FanOut(
				ctx,
				weatherDataChan,
				processWeatherDataFunc,
				workerCount,
			)
			
			// 3. Process results and errors
			var batchResults []WeatherResult
			var batchWarnings []WeatherResult
			var batchErrors []error
			
			// Process all results in this batch
			wg := sync.WaitGroup{}
			wg.Add(2)
			
			// Collect results
			go func() {
				defer wg.Done()
				for result := range resultChan {
					batchResults = append(batchResults, result)
					
					if result.WarningLevel > 0 {
						batchWarnings = append(batchWarnings, result)
					}
					
					mu.Lock()
					processedCount++
					totalTime += result.ProcessingTime
					mu.Unlock()
				}
			}()
			
			// Collect errors
			go func() {
				defer wg.Done()
				for err := range errChan {
					if err != nil {
						batchErrors = append(batchErrors, err)
					}
				}
			}()
			
			// Wait for all results and errors to be collected
			wg.Wait()
			
			// 4. Print batch statistics
			fmt.Printf("\nProcessed %d locations (warnings: %d, errors: %d)\n", 
				len(batchResults), len(batchWarnings), len(batchErrors))
			
			// 5. Print warnings if any
			if len(batchWarnings) > 0 {
				fmt.Println("\n⚠️ Weather Warnings:")
				for _, warning := range batchWarnings {
					level := strings.Repeat("!", warning.WarningLevel)
					fmt.Printf("  %s [%s] %s\n", warning.Location, level, warning.Analysis)
				}
			}
			
			// 6. Print errors if any
			if len(batchErrors) > 0 {
				fmt.Println("\n❌ Errors:")
				for _, err := range batchErrors {
					fmt.Printf("  %s\n", err)
				}
			}
			
			// 7. Print overall statistics
			mu.Lock()
			warningCount += len(batchWarnings)
			if processedCount > 0 {
				avgTime := totalTime / time.Duration(processedCount)
				fmt.Printf("\n📊 Statistics: Processed %d locations, %d warnings, avg processing time: %v\n",
					processedCount, warningCount, avgTime)
			}
			mu.Unlock()
		}
	}
}

// fetchWeatherData simulates fetching weather data from multiple locations
func fetchWeatherData(ctx context.Context, locations []string) <-chan WeatherData {
	out := make(chan WeatherData)
	
	go func() {
		defer close(out)
		
		for _, location := range locations {
			// Check if context is canceled
			select {
			case <-ctx.Done():
				return
			default:
				// Simulate network delay
				time.Sleep(50 * time.Millisecond)
				
				// Generate simulated weather data
				data := WeatherData{
					Location:    location,
					Temperature: 10 + rand.Float64()*30,              // 10°C to 40°C
					WindSpeed:   rand.Float64() * 70,                 // 0 to 70 km/h
					Humidity:    rand.Intn(100),                      // 0% to 100%
					Timestamp:   time.Now(),
				}
				
				// Rarely simulate a fetch error by not sending data
				if rand.Intn(20) != 0 {
					select {
					case <-ctx.Done():
						return
					case out <- data:
						// Data sent successfully
					}
				}
			}
		}
	}()
	
	return out
}

// processWeatherDataFunc analyzes weather data and returns results
// This is the worker function used in the Fan-Out pattern
func processWeatherDataFunc(ctx context.Context, data WeatherData) (WeatherResult, error) {
	// Simulate worker ID (in a real app, this might come from a worker pool)
	workerID := rand.Intn(100)
	
	// Simulate processing time - more extreme conditions take longer to analyze
	processingFactor := (data.WindSpeed / 10) + math.Abs(data.Temperature-20)/5
	processingTime := time.Duration(50+rand.Intn(50)) * time.Millisecond
	processingTime += time.Duration(processingFactor * 10) * time.Millisecond
	
	startTime := time.Now()
	
	// Simulate work being done
	select {
	case <-ctx.Done():
		return WeatherResult{}, ctx.Err()
	case <-time.After(processingTime):
		// Continue processing
	}
	
	// Determine warning level and analysis
	var warningLevel int
	var analysis string
	
	// Extreme temperature check
	if data.Temperature > 35 {
		warningLevel = 3
		analysis = "EXTREME HEAT ALERT! "
	} else if data.Temperature > 30 {
		warningLevel = 2
		analysis = "Heat warning! "
	} else if data.Temperature < 0 {
		warningLevel = 2
		analysis = "Freezing conditions! "
	}
	
	// Wind speed check
	if data.WindSpeed > 60 {
		warningLevel = 3
		analysis += "DANGEROUS WIND CONDITIONS! "
	} else if data.WindSpeed > 40 {
		warningLevel = max(warningLevel, 2)
		analysis += "Strong winds! "
	} else if data.WindSpeed > 20 {
		warningLevel = max(warningLevel, 1)
		analysis += "Moderate winds. "
	}
	
	// Humidity check
	if data.Humidity > 90 {
		warningLevel = max(warningLevel, 1)
		analysis += "Very humid. "
	} else if data.Humidity < 20 {
		warningLevel = max(warningLevel, 1)
		analysis += "Very dry conditions. "
	}
	
	// If no warnings, provide normal analysis
	if warningLevel == 0 {
		analysis = fmt.Sprintf("Normal conditions: %.1f°C, %.1f km/h wind, %d%% humidity", 
			data.Temperature, data.WindSpeed, data.Humidity)
	} else {
		analysis += fmt.Sprintf("Readings: %.1f°C, %.1f km/h wind, %d%% humidity", 
			data.Temperature, data.WindSpeed, data.Humidity)
	}
	
	// Create result
	result := WeatherResult{
		Location:      data.Location,
		Analysis:      analysis,
		WarningLevel:  warningLevel,
		ProcessedBy:   workerID,
		ProcessingTime: time.Since(startTime),
	}
	
	return result, nil
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
