package main

import (
	"fmt"
	"os"
	"time"

	"github.com/edgardnogueira/go-patterns/structural/proxy"
)

func main() {
	fmt.Println("Proxy Pattern Example")
	fmt.Println("=====================")

	// Example 1: Basic Real Image
	example1BasicRealImage()

	// Example 2: Virtual Proxy (Lazy Loading)
	example2VirtualProxy()

	// Example 3: Protection Proxy (Access Control)
	example3ProtectionProxy()

	// Example 4: Caching Proxy
	example4CachingProxy()

	// Example 5: Logging Proxy
	example5LoggingProxy()

	// Example 6: Metrics Proxy
	example6MetricsProxy()

	// Example 7: Remote Proxy
	example7RemoteProxy()

	// Example 8: Proxy Chaining
	example8ProxyChaining()

	// Example 9: Real-world Scenario
	example9RealWorldScenario()
}

func example1BasicRealImage() {
	fmt.Println("\n=== Example 1: Basic Real Image ===")

	// Create a real image
	image, err := proxy.NewRealImage("sample.jpg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get image information
	fmt.Printf("Image: %s (%dx%d, %d bytes)\n",
		image.GetFilename(),
		image.GetWidth(),
		image.GetHeight(),
		image.GetSize())

	// Display the image
	fmt.Println("Displaying the image directly:")
	image.Display()
}

func example2VirtualProxy() {
	fmt.Println("\n=== Example 2: Virtual Proxy (Lazy Loading) ===")

	// Create a virtual proxy (lazy loading)
	virtualProxy := proxy.NewVirtualProxy("large_image.jpg")

	// Note that the image isn't loaded yet
	fmt.Printf("Image: %s (not loaded yet)\n", virtualProxy.GetFilename())
	fmt.Printf("Is image loaded? %v\n", virtualProxy.IsLoaded())

	// Display the image - this will trigger loading
	fmt.Println("Displaying the image (will trigger loading):")
	virtualProxy.Display()

	// Now the image is loaded
	fmt.Printf("Is image loaded now? %v\n", virtualProxy.IsLoaded())
	fmt.Printf("Image dimensions: %dx%d\n", virtualProxy.GetWidth(), virtualProxy.GetHeight())
}

func example3ProtectionProxy() {
	fmt.Println("\n=== Example 3: Protection Proxy (Access Control) ===")

	// Create a real image
	realImage, _ := proxy.NewRealImage("confidential.jpg")

	// Define users with different roles
	adminUser := &proxy.User{Username: "admin", Role: "admin"}
	guestUser := &proxy.User{Username: "guest", Role: "guest"}

	// Create protection proxies for each user
	adminProxy := proxy.NewProtectionProxy(realImage, adminUser)
	guestProxy := proxy.NewProtectionProxy(realImage, guestUser)

	// Admin should have access
	fmt.Println("Admin attempting to access the image:")
	err := adminProxy.Display()
	if err != nil {
		fmt.Printf("Access denied: %v\n", err)
	}

	// Guest should be denied
	fmt.Println("\nGuest attempting to access the image:")
	err = guestProxy.Display()
	if err != nil {
		fmt.Printf("Access denied: %v\n", err)
	}

	// Update allowed roles to include 'guest'
	fmt.Println("\nUpdating access control to allow guests:")
	guestProxy.SetAllowedRoles([]string{"admin", "editor", "guest"})

	// Guest should now have access
	fmt.Println("Guest attempting to access the image again:")
	err = guestProxy.Display()
	if err != nil {
		fmt.Printf("Access denied: %v\n", err)
	}
}

func example4CachingProxy() {
	fmt.Println("\n=== Example 4: Caching Proxy ===")

	// Create a caching proxy with 30 second expiration and max 5 items
	cachingProxy := proxy.NewCachingProxy(30*time.Second, 5)

	// Access an image - first time it will be a cache miss
	fmt.Println("First access (should be a cache miss):")
	cachingProxy.Display("beach.jpg")

	// Access the same image again - should be a cache hit
	fmt.Println("\nSecond access (should be a cache hit):")
	cachingProxy.Display("beach.jpg")

	// Access a different image
	fmt.Println("\nAccessing a different image:")
	cachingProxy.Display("mountain.jpg")

	// Check cache statistics
	stats := cachingProxy.GetCacheStats()
	fmt.Println("\nCache Statistics:")
	fmt.Printf("Cache size: %d/%d\n", stats["size"], stats["max_size"])
	fmt.Printf("Cache expiration: %.0f seconds\n", stats["expiration_ms"].(float64)/1000)
	
	// Print stats for each cached item
	fmt.Println("\nPer-item Statistics:")
	items := stats["items"].(map[string]map[string]interface{})
	for filename, itemStats := range items {
		fmt.Printf("  %s:\n", filename)
		fmt.Printf("    Accesses: %.0f\n", itemStats["access_count"])
		fmt.Printf("    Age: %.1f seconds\n", itemStats["age_seconds"])
		fmt.Printf("    Last accessed: %.1f seconds ago\n", itemStats["last_access_ago"])
	}

	// Clear the cache
	fmt.Println("\nClearing the cache:")
	cachingProxy.ClearCache()
}

func example5LoggingProxy() {
	fmt.Println("\n=== Example 5: Logging Proxy ===")

	// Create a real image
	realImage, _ := proxy.NewRealImage("vacation.jpg")

	// Create a logging proxy
	loggingProxy := proxy.NewLoggingProxy(realImage, proxy.INFO)

	// Perform operations with logging
	fmt.Println("Performing operations with logging:")
	loggingProxy.Display()
	loggingProxy.GetWidth()
	loggingProxy.GetHeight()
	loggingProxy.GetMetadata()

	// Change log level to only show warnings and errors
	fmt.Println("\nChanging log level to WARNING:")
	loggingProxy.SetLogLevel(proxy.WARNING)
	loggingProxy.Display() // This shouldn't generate logs at INFO level
}

func example6MetricsProxy() {
	fmt.Println("\n=== Example 6: Metrics Proxy ===")

	// Create a real image
	realImage, _ := proxy.NewRealImage("chart.jpg")

	// Create a metrics proxy
	metricsProxy := proxy.NewMetricsProxy(realImage)

	// Set up a metrics hook
	metricsProxy.SetMetricsHook(func(operation string, metricType proxy.MetricType, value float64) {
		// In a real application, this could send metrics to a monitoring system
		fmt.Printf("Metric: %s - %s = %.2f\n", operation, metricType, value)
	})

	// Perform operations to collect metrics
	fmt.Println("Performing operations to collect metrics:")
	for i := 0; i < 3; i++ {
		metricsProxy.Display()
		metricsProxy.GetWidth()
		metricsProxy.GetHeight()
		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)
	}

	// Print a metrics report
	fmt.Println("\nMetrics Report:")
	metricsProxy.PrintMetricsReport()
}

func example7RemoteProxy() {
	fmt.Println("\n=== Example 7: Remote Proxy ===")

	// Create a remote proxy
	remoteProxy := proxy.NewRemoteProxy("https://images.example.com/api", "remote_scenery.jpg")

	// Set a short cache duration for demonstration
	remoteProxy.SetCacheDuration(10 * time.Second)

	// Access the remote image
	fmt.Println("Accessing remote image (initial fetch):")
	remoteProxy.Display()

	// Get image information
	fmt.Printf("\nRemote image: %s (%dx%d, %d bytes)\n",
		remoteProxy.GetFilename(),
		remoteProxy.GetWidth(),
		remoteProxy.GetHeight(),
		remoteProxy.GetSize())

	// Access again - should use cache
	fmt.Println("\nAccessing remote image again (using cache):")
	remoteProxy.Display()

	// Check cache status
	cacheStatus := remoteProxy.GetCacheStatus()
	fmt.Println("\nCache Status:")
	for key, value := range cacheStatus {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Clear cache and access again
	fmt.Println("\nClearing cache and accessing again:")
	remoteProxy.ClearCache()
	remoteProxy.Display()
}

func example8ProxyChaining() {
	fmt.Println("\n=== Example 8: Proxy Chaining ===")

	// Create a real image
	realImage, _ := proxy.NewRealImage("landscape.jpg")

	// Create an admin user
	adminUser := &proxy.User{Username: "admin", Role: "admin"}

	// Build a chain of proxies
	fmt.Println("Building a chain with multiple proxies:")
	chainedProxy := proxy.NewProxyChain(realImage).
		AddLogging(proxy.INFO).
		AddMetrics().
		AddProtection(adminUser).
		Build()

	// Use the chained proxy
	fmt.Println("\nUsing the chained proxy:")
	chainedProxy.Display()

	// Create a chain with virtual loading and caching
	fmt.Println("\nCreating a chain with virtual loading and caching:")
	optimizedChain := proxy.ProxyPresets{}.NewOptimized("heavy_image.jpg")

	// Use the optimized chain
	fmt.Println("\nUsing the optimized chain:")
	optimizedChain.Display()
}

func example9RealWorldScenario() {
	fmt.Println("\n=== Example 9: Real-world Scenario ===")
	fmt.Println("Image Processing Service with Security, Caching, and Monitoring")

	// Create different user types
	adminUser := &proxy.User{Username: "admin", Role: "admin"}
	editorUser := &proxy.User{Username: "editor", Role: "editor"}
	guestUser := &proxy.User{Username: "guest", Role: "guest"}

	// Set up a common image processing function
	processImage := func(user *proxy.User, imageName string) {
		fmt.Printf("\nUser %s (%s) requesting image: %s\n", user.Username, user.Role, imageName)

		// Create the base image with virtual loading
		baseImage := proxy.NewVirtualProxy(imageName)

		// Add protection
		protectedImage := proxy.NewProtectionProxy(baseImage, user)
		protectedImage.SetAllowedRoles([]string{"admin", "editor"}) // Only admin and editor can access

		// Add logging (to a file in a real application)
		loggingImage := proxy.NewLoggingProxy(protectedImage, proxy.INFO)

		// Add metrics
		monitoredImage := proxy.NewMetricsProxy(loggingImage)

		// Try to access the image
		err := monitoredImage.Display()
		if err != nil {
			fmt.Printf("Result: %v\n", err)
		} else {
			fmt.Printf("Result: Successfully processed image %s (%dx%d)\n",
				imageName, monitoredImage.GetWidth(), monitoredImage.GetHeight())
		}
	}

	// Process images with different users
	processImage(adminUser, "financial_report.jpg")
	processImage(editorUser, "marketing_campaign.jpg")
	processImage(guestUser, "confidential_data.jpg") // Should be denied

	// Using a caching proxy for a gallery application
	fmt.Println("\nGallery Application with Caching:")
	
	// Create a caching proxy
	gallery := proxy.NewCachingProxy(5*time.Minute, 100)
	
	// Simulate browsing a gallery
	images := []string{"gallery/image1.jpg", "gallery/image2.jpg", "gallery/image1.jpg", "gallery/image3.jpg", "gallery/image1.jpg"}
	
	for _, img := range images {
		fmt.Printf("Accessing %s: ", img)
		gallery.Display(img)
	}
	
	// Print cache statistics
	stats := gallery.GetCacheStats()
	fmt.Println("\nGallery Cache Statistics:")
	fmt.Printf("Images in cache: %d\n", stats["size"])
	
	// Print the most accessed images
	items := stats["items"].(map[string]map[string]interface{})
	fmt.Println("Most accessed images:")
	for img, itemStats := range items {
		fmt.Printf("  %s: %.0f views\n", img, itemStats["access_count"])
	}
}
