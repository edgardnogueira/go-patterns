package singleton

import (
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Error("Instances are not the same")
	}

	instance1.IncrementCount()
	if instance1.GetCount() != 1 {
		t.Errorf("Expected count to be 1, got %d", instance1.GetCount())
	}

	if instance2.GetCount() != 1 {
		t.Errorf("Expected count to be 1, got %d", instance2.GetCount())
	}
}

func TestConcurrentSingleton(t *testing.T) {
	const numGoroutines = 100
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			instance := GetInstance()
			instance.IncrementCount()
			wg.Done()
		}()
	}

	wg.Wait()
	instance := GetInstance()
	if instance.GetCount() != numGoroutines {
		t.Errorf("Expected count to be %d, got %d", numGoroutines, instance.GetCount())
	}
}
