// Package di_container demonstrates a lightweight dependency injection container in Go.
package di_container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container is a dependency injection container
type Container struct {
	mu           sync.RWMutex
	services     map[reflect.Type]interface{}
	factories    map[reflect.Type]interface{}
	singletons   map[reflect.Type]interface{}
	constructors map[reflect.Type]interface{}
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	return &Container{
		services:     make(map[reflect.Type]interface{}),
		factories:    make(map[reflect.Type]interface{}),
		singletons:   make(map[reflect.Type]interface{}),
		constructors: make(map[reflect.Type]interface{}),
	}
}

// Register adds a service instance to the container
func (c *Container) Register(service interface{}) {
	if service == nil {
		return
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	t := reflect.TypeOf(service)
	c.services[t] = service
}

// RegisterType registers a service implementation for an interface
func (c *Container) RegisterType(interfacePtr, implementation interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make sure interfacePtr is a pointer to an interface
	interfacePtrType := reflect.TypeOf(interfacePtr)
	if interfacePtrType.Kind() != reflect.Ptr || interfacePtrType.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("first argument must be a pointer to an interface")
	}
	
	// Get the actual interface type
	interfaceType := interfacePtrType.Elem()
	
	// Make sure implementation implements the interface
	implType := reflect.TypeOf(implementation)
	if !implType.Implements(interfaceType) {
		return fmt.Errorf("implementation does not implement the interface %v", interfaceType)
	}
	
	c.services[interfaceType] = implementation
	return nil
}

// RegisterFactory registers a factory function for a service
// The factory will be called each time the service is requested
func (c *Container) RegisterFactory(interfacePtr interface{}, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make sure interfacePtr is a pointer to an interface
	interfacePtrType := reflect.TypeOf(interfacePtr)
	if interfacePtrType.Kind() != reflect.Ptr || interfacePtrType.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("first argument must be a pointer to an interface")
	}
	
	// Get the actual interface type
	interfaceType := interfacePtrType.Elem()
	
	// Make sure factory is a function
	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return fmt.Errorf("factory must be a function")
	}
	
	// Make sure factory returns the interface
	if factoryType.NumOut() < 1 || !factoryType.Out(0).Implements(interfaceType) {
		return fmt.Errorf("factory function must return an implementation of the interface %v", interfaceType)
	}
	
	c.factories[interfaceType] = factory
	return nil
}

// RegisterSingleton registers a factory function that will be called only once
// to create a singleton instance of the service
func (c *Container) RegisterSingleton(interfacePtr interface{}, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make sure interfacePtr is a pointer to an interface
	interfacePtrType := reflect.TypeOf(interfacePtr)
	if interfacePtrType.Kind() != reflect.Ptr || interfacePtrType.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("first argument must be a pointer to an interface")
	}
	
	// Get the actual interface type
	interfaceType := interfacePtrType.Elem()
	
	// Make sure factory is a function
	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return fmt.Errorf("factory must be a function")
	}
	
	// Make sure factory returns the interface
	if factoryType.NumOut() < 1 || !factoryType.Out(0).Implements(interfaceType) {
		return fmt.Errorf("factory function must return an implementation of the interface %v", interfaceType)
	}
	
	c.constructors[interfaceType] = factory
	return nil
}

// RegisterInstance registers a concrete instance for an interface
func (c *Container) RegisterInstance(interfacePtr interface{}, instance interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Make sure interfacePtr is a pointer to an interface
	interfacePtrType := reflect.TypeOf(interfacePtr)
	if interfacePtrType.Kind() != reflect.Ptr || interfacePtrType.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("first argument must be a pointer to an interface")
	}
	
	// Get the actual interface type
	interfaceType := interfacePtrType.Elem()
	
	// Make sure instance implements the interface
	instanceType := reflect.TypeOf(instance)
	if !instanceType.Implements(interfaceType) {
		return fmt.Errorf("instance does not implement the interface %v", interfaceType)
	}
	
	c.services[interfaceType] = instance
	return nil
}

// Resolve retrieves a service from the container
func (c *Container) Resolve(interfacePtr interface{}) error {
	if interfacePtr == nil {
		return fmt.Errorf("cannot resolve nil interface")
	}
	
	// Make sure interfacePtr is a pointer to an interface
	targetType := reflect.TypeOf(interfacePtr)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	
	// Get the actual interface type
	elemType := targetType.Elem()
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// First check if we have a singleton instance
	if singleton, ok := c.singletons[elemType]; ok {
		reflect.ValueOf(interfacePtr).Elem().Set(reflect.ValueOf(singleton))
		return nil
	}
	
	// Check if we have a registered service
	if service, ok := c.services[elemType]; ok {
		reflect.ValueOf(interfacePtr).Elem().Set(reflect.ValueOf(service))
		return nil
	}
	
	// Check if we have a factory
	if factory, ok := c.factories[elemType]; ok {
		result := reflect.ValueOf(factory).Call(nil)
		if len(result) > 0 && !result[0].IsNil() {
			reflect.ValueOf(interfacePtr).Elem().Set(result[0])
			return nil
		}
		return fmt.Errorf("factory returned nil or invalid result")
	}
	
	// Check if we have a singleton constructor
	if constructor, ok := c.constructors[elemType]; ok {
		result := reflect.ValueOf(constructor).Call(nil)
		if len(result) > 0 && !result[0].IsNil() {
			instance := result[0].Interface()
			// Store the instance for future resolves
			c.singletons[elemType] = instance
			reflect.ValueOf(interfacePtr).Elem().Set(reflect.ValueOf(instance))
			return nil
		}
		return fmt.Errorf("singleton constructor returned nil or invalid result")
	}
	
	return fmt.Errorf("no registration found for %v", elemType)
}

// AutoWire injects dependencies into the struct's fields 
// based on their types
func (c *Container) AutoWire(target interface{}) error {
	if target == nil {
		return fmt.Errorf("cannot autowire nil target")
	}
	
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer to a struct")
	}
	
	targetElem := targetValue.Elem()
	if targetElem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	
	// Process each field in the struct
	targetType := targetElem.Type()
	for i := 0; i < targetElem.NumField(); i++ {
		field := targetElem.Field(i)
		
		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		
		// Skip non-interface fields
		fieldType := targetType.Field(i).Type
		if fieldType.Kind() != reflect.Interface {
			continue
		}
		
		// Try to resolve the dependency
		c.mu.RLock()
		
		// Check for singleton
		if singleton, ok := c.singletons[fieldType]; ok {
			field.Set(reflect.ValueOf(singleton))
			c.mu.RUnlock()
			continue
		}
		
		// Check for service
		if service, ok := c.services[fieldType]; ok {
			field.Set(reflect.ValueOf(service))
			c.mu.RUnlock()
			continue
		}
		
		// Check for factory
		if factory, ok := c.factories[fieldType]; ok {
			c.mu.RUnlock()
			result := reflect.ValueOf(factory).Call(nil)
			if len(result) > 0 && !result[0].IsNil() {
				field.Set(result[0])
			}
			continue
		}
		
		// Check for singleton constructor
		if constructor, ok := c.constructors[fieldType]; ok {
			result := reflect.ValueOf(constructor).Call(nil)
			if len(result) > 0 && !result[0].IsNil() {
				instance := result[0].Interface()
				// Store the instance for future resolves
				c.singletons[fieldType] = instance
				field.Set(reflect.ValueOf(instance))
			}
			c.mu.RUnlock()
			continue
		}
		
		c.mu.RUnlock()
		// No registration found, but we continue - field stays as its zero value
	}
	
	return nil
}

// Clear removes all registrations from the container
func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.services = make(map[reflect.Type]interface{})
	c.factories = make(map[reflect.Type]interface{})
	c.singletons = make(map[reflect.Type]interface{})
	c.constructors = make(map[reflect.Type]interface{})
}

//
// Example interfaces and implementations
//

// Logger is a simple logging interface
type Logger interface {
	Log(message string)
}

// ConsoleLogger is a basic implementation of Logger
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

// FileLogger logs to a file
type FileLogger struct {
	FilePath string
}

func (l *FileLogger) Log(message string) {
	fmt.Printf("Writing to file %s: %s\n", l.FilePath, message)
}

// LoggerFactory creates a new logger
func LoggerFactory() Logger {
	return &ConsoleLogger{}
}

// Repository defines data access methods
type Repository interface {
	FindByID(id string) (interface{}, error)
	Save(id string, data interface{}) error
}

// MemoryRepository implements Repository with in-memory storage
type MemoryRepository struct {
	data map[string]interface{}
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]interface{}),
	}
}

func (r *MemoryRepository) FindByID(id string) (interface{}, error) {
	if data, ok := r.data[id]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("item not found: %s", id)
}

func (r *MemoryRepository) Save(id string, data interface{}) error {
	r.data[id] = data
	return nil
}

// UserService uses DI container
type UserService struct {
	Logger Logger           // Will be injected
	Repo   Repository       // Will be injected
	Config *ServiceConfig   // Will be injected
}

// ServiceConfig holds configuration
type ServiceConfig struct {
	Timeout int
	BaseURL string
}

// NewUserService creates a new UserService with dependencies injected
func NewUserService(logger Logger, repo Repository, config *ServiceConfig) *UserService {
	return &UserService{
		Logger: logger,
		Repo:   repo,
		Config: config,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(id, name string) error {
	s.Logger.Log(fmt.Sprintf("Creating user: %s", name))
	
	user := struct {
		ID   string
		Name string
	}{
		ID:   id,
		Name: name,
	}
	
	if err := s.Repo.Save(id, user); err != nil {
		s.Logger.Log(fmt.Sprintf("Error saving user: %v", err))
		return err
	}
	
	s.Logger.Log(fmt.Sprintf("User created with timeout: %d", s.Config.Timeout))
	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (interface{}, error) {
	s.Logger.Log(fmt.Sprintf("Getting user: %s", id))
	return s.Repo.FindByID(id)
}
