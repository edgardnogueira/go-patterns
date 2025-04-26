package main

import (
	"testing"
)

// Test the fat Worker interface
func TestWorkerInterface(t *testing.T) {
	// Create workers of different types
	manager := &Manager{Name: "TestManager"}
	developer := &Developer{Name: "TestDev"}
	intern := &Intern{Name: "TestIntern"}
	
	// We can create a slice of Worker with all types
	workers := []Worker{manager, developer, intern}
	
	// Test core worker functions that all workers should implement properly
	t.Run("core worker functions", func(t *testing.T) {
		for _, worker := range workers {
			// These should work for all worker types
			worker.DoWork()
			worker.TakeBreak()
			worker.GetPaid()
		}
	})
	
	// Test management functions that only make sense for managers
	t.Run("management functions", func(t *testing.T) {
		for _, worker := range workers {
			// These methods should only make sense for Manager,
			// but Developer and Intern must implement them anyway
			worker.ManageOthers()
			worker.GenerateReports()
		}
		
		// Specifically check that non-managers have empty implementations
		developer.ManageOthers() // Empty implementation
		intern.ManageOthers()    // Empty implementation
		
		// This demonstrates ISP violation: Developer and Intern are forced
		// to implement methods that don't make sense for them
	})
}

// Test Worker implementations for inappropriate method implementations
func TestFatInterfaceProblems(t *testing.T) {
	developer := &Developer{Name: "TestDev"}
	intern := &Intern{Name: "TestIntern"}
	
	t.Run("developer shouldn't manage but must implement", func(t *testing.T) {
		// Developer must implement ManageOthers even though it's not appropriate
		developer.ManageOthers()
		
		// Instead of checking specific behavior, we're just documenting the issue:
		// - The method exists but does nothing useful
		// - This is a violation of ISP
	})
	
	t.Run("intern shouldn't generate reports but must implement", func(t *testing.T) {
		// Intern must implement GenerateReports even though it's not appropriate
		intern.GenerateReports()
		
		// Again, just documenting the issue:
		// - The method exists but does nothing useful
		// - This is a violation of ISP
	})
}

// MockWorkerFatInterface is a test implementation of the fat Worker interface
type MockWorkerFatInterface struct {
	// We need fields for all methods, even ones we don't care about
	DoWorkCalled         bool
	TakeBreakCalled      bool
	GetPaidCalled        bool
	ManageOthersCalled   bool  // Don't need this
	GenerateReportsCalled bool // Don't need this
	AttendMeetingsCalled bool  // Don't need this
	WorkOvertimeCalled   bool  // Don't need this
}

func (m *MockWorkerFatInterface) DoWork() {
	m.DoWorkCalled = true
}

func (m *MockWorkerFatInterface) TakeBreak() {
	m.TakeBreakCalled = true
}

func (m *MockWorkerFatInterface) GetPaid() {
	m.GetPaidCalled = true
}

// Must implement these methods even though we don't need them
func (m *MockWorkerFatInterface) ManageOthers() {
	m.ManageOthersCalled = true
}

func (m *MockWorkerFatInterface) GenerateReports() {
	m.GenerateReportsCalled = true
}

func (m *MockWorkerFatInterface) AttendMeetings() {
	m.AttendMeetingsCalled = true
}

func (m *MockWorkerFatInterface) WorkOvertime() {
	m.WorkOvertimeCalled = true
}

// Test using a mock implementation of fat interface
func TestMockWorkerFatInterface(t *testing.T) {
	t.Run("mock worker must implement all methods", func(t *testing.T) {
		// Create a mock worker
		mockWorker := &MockWorkerFatInterface{}
		
		// Even though we only want to test DoWork functionality,
		// we are forced to implement all other methods
		var worker Worker = mockWorker
		
		// Use the worker for just the one method we care about
		worker.DoWork()
		
		// Verify only the method we care about was called
		if !mockWorker.DoWorkCalled {
			t.Error("DoWork was not called")
		}
		
		// The issue: We had to implement all these other methods
		// just to test the one method we care about
		t.Log("Had to implement 7 methods to test 1 method")
	})
}

// Test the ISP violation impact on a worker client
func TestWorkerClient(t *testing.T) {
	t.Run("client code affected by fat interface", func(t *testing.T) {
		// Create some workers
		manager := &Manager{Name: "Alice"}
		developer := &Developer{Name: "Bob"}
		intern := &Intern{Name: "Charlie"}
		
		workers := []Worker{manager, developer, intern}
		
		// Using our client code
		// This client code must call methods on all workers that don't make
		// sense for some of them, violating client expectations
		
		// Morning routine should be ok for all
		for _, worker := range workers {
			worker.DoWork()
		}
		
		// But management activities shouldn't apply to all
		// Yet our interface forces us to call these methods on all workers
		for _, worker := range workers {
			worker.ManageOthers()  // Only managers should do this
			worker.GenerateReports() // Only managers should do this
		}
		
		// This results in empty or inappropriate method calls
	})
}

func TestISPViolationIssues(t *testing.T) {
	t.Run("highlight issues with ISP violations", func(t *testing.T) {
		// This "test" highlights the issues with ISP violations
		
		// ISSUE 1: Forced implementation of irrelevant methods
		t.Log("ISSUE 1: Forced implementation of irrelevant methods")
		t.Log("- Developer and Intern must implement ManageOthers and GenerateReports")
		t.Log("- Mock implementations must implement all methods to test just one")
		t.Log("- This leads to empty or inappropriate implementations")
		
		// ISSUE 2: Interface bloat
		t.Log("ISSUE 2: Interface bloat")
		t.Log("- Worker interface has 7 methods when many clients need only 2-3")
		t.Log("- Adding new methods to Worker requires updating all implementations")
		t.Log("- This makes the interface increasingly difficult to implement")
		
		// ISSUE 3: Client violation of expectations
		t.Log("ISSUE 3: Client violation of expectations")
		t.Log("- Clients expect methods to do something meaningful")
		t.Log("- But many methods do nothing or return errors for certain types")
		t.Log("- This leads to confusing behavior and potential bugs")
		
		// ISSUE 4: Testing complexity
		t.Log("ISSUE 4: Testing complexity")
		t.Log("- Testing one role requires implementing all methods")
		t.Log("- Mock objects become unnecessarily complex")
		t.Log("- More difficult to isolate specific behaviors")
		
		// This is not a real test, just documentation
		t.Skip("This is not a real test, just documentation of ISP violation issues")
	})
}
