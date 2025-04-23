package singleton

import (
	"sync"
)

type singleton struct {
	count int
}

var (
	instance *singleton
	once     sync.Once
)

// GetInstance returns the single instance of singleton
func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{count: 0}
	})
	return instance
}

// IncrementCount increments the counter
func (s *singleton) IncrementCount() {
	s.count++
}

// GetCount returns the current count
func (s *singleton) GetCount() int {
	return s.count
}
