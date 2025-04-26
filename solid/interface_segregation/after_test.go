package main

import (
	"testing"
)

// TestBasicWorker tests the BasicWorker interface
func TestBasicWorker(t *testing.T) {
	// Create workers of different types
	manager := &FullManager{Name: "TestManager"}
	developer := &FullDeveloper{Name: "TestDev"}
	intern := &FullIntern{Name: "TestIntern"}
	
	// We can create a slice of BasicWorker with all types
	workers := []BasicWorker{manager, developer, intern}
	
	// Test that all workers implement DoWork correctly
	t.Run("all workers can do work", func(t *testing.T) {
		for _, worker := range workers {
			// This just tests that the method doesn't panic
			worker.DoWork()
		}
	})
	
	// Test that all workers implement TakeBreak correctly
	t.Run("all workers can take breaks", func(t *testing.T) {
		for _, worker := range workers {
			// This just tests that the method doesn't panic
			worker.TakeBreak()
		}
	})
	
	// Test that all workers implement GetPaid correctly
	t.Run("all workers can get paid", func(t *testing.T) {
		for _, worker := range workers {
			// This just tests that the method doesn't panic
			worker.GetPaid()
		}
	})
	
	// Test the MorningRoutine function
	t.Run("morning routine works for all workers", func(t *testing.T) {
		// We can pass any BasicWorker to this function
		MorningRoutine(workers)
	})
}

// TestTeamMember tests the TeamMember interface
func TestTeamMember(t *testing.T) {
	// Create team members of different types
	manager := &FullManager{Name: "TestManager"}
	developer := &FullDeveloper{Name: "TestDev"}
	intern := &FullIntern{Name: "TestIntern"}
	
	// We can create a slice of TeamMember with all types
	members := []TeamMember{manager, developer, intern}
	
	// Test that all members implement AttendMeetings correctly
	t.Run("all team members can attend meetings", func(t *testing.T) {
		for _, member := range members {
			// This just tests that the method doesn't panic
			member.AttendMeetings()
		}
	})
	
	// Test the TeamMeeting function
	t.Run("team meeting works for all team members", func(t *testing.T) {
		// We can pass any TeamMember to this function
		TeamMeeting(members)
	})
}

// TestOvertimeWorker tests the OvertimeWorker interface
func TestOvertimeWorker(t *testing.T) {
	// Create overtime workers of different types
	manager := &FullManager{Name: "TestManager"}
	developer := &FullDeveloper{Name: "TestDev"}
	
	// We can create a slice of OvertimeWorker with only the types that implement it
	workers := []OvertimeWorker{manager, developer}
	
	// Test that all overtime workers implement WorkOvertime correctly
	t.Run("overtime workers can work overtime", func(t *testing.T) {
		for _, worker := range workers {
			// This just tests that the method doesn't panic
			worker.WorkOvertime()
		}
	})
	
	// Test the OvertimeWork function
	t.Run("overtime work function works for overtime workers", func(t *testing.T) {
		// We can pass any OvertimeWorker to this function
		OvertimeWork(workers)
	})
	
	// Compile-time verification: intern doesn't implement OvertimeWorker
	// The following would cause a compile error:
	// intern := &FullIntern{Name: "TestIntern"}
	// workers = append(workers, intern)
}

// TestManagerRole tests the ManagerRole interface
func TestManagerRole(t *testing.T) {
	// Create managers
	manager := &FullManager{Name: "TestManager"}
	
	// We can create a slice of ManagerRole with only the types that implement it
	managers := []ManagerRole{manager}
	
	// Test that managers implement ManageOthers correctly
	t.Run("managers can manage others", func(t *testing.T) {
		for _, mgr := range managers {
			// This just tests that the method doesn't panic
			mgr.ManageOthers()
		}
	})
	
	// Test that managers implement GenerateReports correctly
	t.Run("managers can generate reports", func(t *testing.T) {
		for _, mgr := range managers {
			// This just tests that the method doesn't panic
			mgr.GenerateReports()
		}
	})
	
	// Test the ManagementActivities function
	t.Run("management activities function works for managers", func(t *testing.T) {
		// We can pass any ManagerRole to this function
		ManagementActivities(managers)
	})
	
	// Compile-time verification: developer and intern don't implement ManagerRole
	// The following would cause compile errors:
	// developer := &FullDeveloper{Name: "TestDev"}
	// managers = append(managers, developer)
	// intern := &FullIntern{Name: "TestIntern"}
	// managers = append(managers, intern)
}

// TestISPComposite tests how the different interfaces work together
func TestISPComposite(t *testing.T) {
	t.Run("interfaces can be composed as needed", func(t *testing.T) {
		// Create workers of different types
		manager := &FullManager{Name: "TestManager"}
		developer := &FullDeveloper{Name: "TestDev"}
		intern := &FullIntern{Name: "TestIntern"}
		
		// We can test different capability combinations
		
		// These functions expect different interfaces
		MorningRoutine([]BasicWorker{manager, developer, intern})
		TeamMeeting([]TeamMember{manager, developer, intern})
		OvertimeWork([]OvertimeWorker{manager, developer})
		ManagementActivities([]ManagerRole{manager})
		
		// We can also compose functions that use different interfaces
		workers := []BasicWorker{manager, developer, intern}
		MorningRoutine(workers)
		
		overtimeCapable := []BasicWorker{manager, developer}
		MorningRoutine(overtimeCapable)
		
		// The ISP allows us to only require the capabilities we actually need
	})
}

// Create a mock worker for further testing
type MockWorker struct {
	DoWorkCalled     bool
	TakeBreakCalled  bool
	GetPaidCalled    bool
}

func (m *MockWorker) DoWork() {
	m.DoWorkCalled = true
}

func (m *MockWorker) TakeBreak() {
	m.TakeBreakCalled = true
}

func (m *MockWorker) GetPaid() {
	m.GetPaidCalled = true
}

// Test using a mock implementation
func TestMockWorker(t *testing.T) {
	t.Run("mock worker implements only needed interface", func(t *testing.T) {
		// Create a mock worker
		mockWorker := &MockWorker{}
		
		// We can use it as a BasicWorker
		workers := []BasicWorker{mockWorker}
		MorningRoutine(workers)
		
		// Verify the methods were called
		if !mockWorker.DoWorkCalled {
			t.Error("DoWork was not called")
		}
		
		// Reset and test MiddayRoutine
		mockWorker = &MockWorker{}
		MiddayRoutine([]BasicWorker{mockWorker})
		
		// Verify the method was called
		if !mockWorker.TakeBreakCalled {
			t.Error("TakeBreak was not called")
		}
		
		// Reset and test EndOfDayRoutine
		mockWorker = &MockWorker{}
		EndOfDayRoutine([]BasicWorker{mockWorker})
		
		// Verify the method was called
		if !mockWorker.GetPaidCalled {
			t.Error("GetPaid was not called")
		}
		
		// The benefit: We only had to implement the methods we actually use
		// We didn't need to implement AttendMeetings, WorkOvertime, etc.
	})
}
