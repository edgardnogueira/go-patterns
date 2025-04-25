# Proxy Pattern

## Intent
The Proxy pattern provides a surrogate or placeholder for another object to control access to it. It creates a representative object that controls access to another object, which may be remote, expensive to create, or require additional security.

## Problem Solved
The Proxy pattern solves the following problems:
- **Controlled Access**: You need to control access to an object, adding security checks or permissions.
- **Lazy Initialization**: You want to postpone the creation of an expensive object until it's actually needed.
- **Caching**: You want to cache results of expensive operations to improve performance.
- **Remote Access**: You need to interact with objects that reside in a different address space.
- **Logging/Monitoring**: You want to add additional behaviors like logging or metrics collection without modifying the original object.

## Structure
![Proxy Pattern Structure](https://github.com/edgardnogueira/go-patterns/raw/main/assets/images/proxy.png)

### Components
1. **Subject**: An interface that defines common operations for both the RealSubject and the Proxy.
2. **RealSubject**: The real object that the proxy represents and controls access to.
3. **Proxy**: Maintains a reference to the RealSubject, controls access to it, and implements the same interface.

### Types of Proxies
- **Virtual Proxy**: Delays the creation and initialization of expensive objects until needed.
- **Protection Proxy**: Controls access to the original object based on access rights.
- **Remote Proxy**: Provides a local representation of an object that resides in a different address space.
- **Caching Proxy**: Stores the results of expensive operations for reuse and improved performance.
- **Logging Proxy**: Adds logging capabilities around method invocations of the original object.
- **Metrics Proxy**: Collects performance metrics and usage statistics for monitoring purposes.

## Implementation in Go

### Subject Interface (Image)
```go
type Image interface {
    Display() error
    GetFilename() string
    GetWidth() int
    GetHeight() int
    GetSize() int64
    GetMetadata() map[string]string
}
```

### Real Subject (RealImage)
```go
type RealImage struct {
    filename string
    width    int
    height   int
    size     int64
    metadata map[string]string
    data     []byte
    loaded   bool
}

func (r *RealImage) Display() error {
    if !r.loaded {
        err := r.loadFromDisk()
        if err != nil {
            return err
        }
    }
    
    fmt.Printf("Displaying image: %s [%dx%d]\n", r.filename, r.width, r.height)
    return nil
}
```

### Base Proxy
```go
type BaseProxy struct {
    realImage Image
}

func (p *BaseProxy) GetFilename() string {
    return p.realImage.GetFilename()
}

// Other forwarding methods...
```

### Virtual Proxy (Lazy Loading)
```go
type VirtualProxy struct {
    filename  string
    realImage Image
    once      sync.Once
    mu        sync.Mutex
}

func (p *VirtualProxy) Display() error {
    err := p.lazyInit()
    if err != nil {
        return fmt.Errorf("virtual proxy initialization error: %w", err)
    }
    
    fmt.Println("Virtual proxy delegating display call to real image")
    return p.realImage.Display()
}
```

### Other Proxy Implementations
- **Protection Proxy**: Adds access control based on user permissions
- **Caching Proxy**: Caches images for faster repeated access
- **Logging Proxy**: Adds logging to method calls
- **Metrics Proxy**: Collects performance metrics
- **Remote Proxy**: Provides a local representation of a remote resource

## When to Use
Use the Proxy pattern when:
1. You need controlled access to an object
2. You want to lazy-load expensive objects
3. You need to add behaviors to objects without changing their code
4. You need to interact with remote objects as if they were local
5. You need to cache results of expensive operations

## Advantages
- Introduces a level of indirection when accessing an object
- Allows open/closed principle by adding functionality without modifying existing code
- Controls the lifecycle of the service object
- Works even when the service object isn't ready or available
- Facilitates separation of concerns and single responsibility principle

## Disadvantages
- Adds another layer of indirection which might impact performance
- Makes the code more complex as it introduces additional classes
- Response from the service might be delayed due to extra layers

## Example Use Case
In our implementation, we've built an image loading system with various proxy types:

```go
// Create a real image
realImage, _ := NewRealImage("landscape.jpg")

// Create an admin user
adminUser := &User{Username: "admin", Role: "admin"}

// Build a chain of proxies
chainedProxy := NewProxyChain(realImage).
    AddLogging(INFO).
    AddMetrics().
    AddProtection(adminUser).
    Build()

// Use the chained proxy
chainedProxy.Display()
```

### Running the Example
A complete example application is provided in the `example` directory. To run it:

```bash
cd structural/proxy/example
go run main.go
```

## Related Patterns
- **Adapter**: While Proxy provides the same interface as its service object, Adapter provides a different interface.
- **Decorator**: Decorator adds responsibilities to an object while Proxy controls access to it.
- **Facade**: Facade provides a simplified interface to a set of subsystems, while Proxy controls access to a single object.

## Go-Specific Considerations
- Go doesn't have inheritance, so the Proxy pattern in Go often uses composition.
- Go interfaces make implementing the Proxy pattern more straightforward.
- Concurrency should be considered when implementing proxies in Go (use sync primitives when needed).
- Go's first-class function support allows for creating simple function-based proxies for some use cases.

## Proxy Chaining
One of the powerful features of this implementation is support for proxy chaining, allowing multiple proxies to be combined:

```go
// Create a chain with multiple proxies
chain := NewProxyChain(realImage).
    AddVirtual().      // Lazy loading
    AddLogging(INFO).  // Logging
    AddMetrics().      // Performance metrics
    AddProtection(user).  // Access control
    Build()

// Access through the chain
chain.Display()
```

This approach lets you compose various proxy behaviors to create complex access control and processing pipelines.

## References
- "Design Patterns: Elements of Reusable Object-Oriented Software" by Gamma, Helm, Johnson, and Vlissides
- "Head First Design Patterns" by Eric Freeman and Elisabeth Robson
- "Design Patterns in Go" by Mario Castro Contreras
