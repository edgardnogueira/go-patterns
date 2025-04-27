package cqrs

import (
	"fmt"
	"sync"
)

// QueryHandler is a function that handles a query
type QueryHandler func(interface{}) (interface{}, error)

// QueryBus routes queries to their appropriate handlers
type QueryBus struct {
	handlers map[string]QueryHandler
	mutex    sync.RWMutex
}

// NewQueryBus creates a new QueryBus
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

// Register registers a query handler for a specific query type
func (b *QueryBus) Register(queryType string, handler QueryHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[queryType] = handler
}

// Dispatch sends a query to its registered handler
func (b *QueryBus) Dispatch(queryType string, query interface{}) (interface{}, error) {
	b.mutex.RLock()
	handler, exists := b.handlers[queryType]
	b.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %s", queryType)
	}

	return handler(query)
}
