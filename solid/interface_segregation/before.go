package main

import (
	"fmt"
)

// Worker is a "fat" interface that forces implementations to provide all methods
// This violates the Interface Segregation Principle by forcing clients to implement
// methods they don't need
type Worker interface {
	DoWork()
	TakeBreak()
	GetPaid()
	ManageOthers()
	GenerateReports()
	AttendMeetings()
	WorkOvertime()
}

// Manager implements Worker but many methods are highly specific to this role
type Manager struct {
	Name string
}

func (m *Manager) DoWork() {
	fmt.Printf("%s is delegating tasks\n", m.Name)
}

func (m *Manager) TakeBreak() {
	fmt.Printf("%s is taking a coffee break\n", m.Name)
}

func (m *Manager) GetPaid() {
	fmt.Printf("%s is getting paid a large salary\n", m.Name)
}

func (m *Manager) ManageOthers() {
	fmt.Printf("%s is managing the team\n", m.Name)
}

func (m *Manager) GenerateReports() {
	fmt.Printf("%s is generating performance reports\n", m.Name)
}

func (m *Manager) AttendMeetings() {
	fmt.Printf("%s is attending a management meeting\n", m.Name)
}

func (m *Manager) WorkOvertime() {
	fmt.Printf("%s is planning tomorrow's work\n", m.Name)
}

// Developer implements Worker but is forced to implement management methods that aren't relevant
type Developer struct {
	Name string
}

func (d *Developer) DoWork() {
	fmt.Printf("%s is writing code\n", d.Name)
}

func (d *Developer) TakeBreak() {
	fmt.Printf("%s is browsing tech news\n", d.Name)
}

func (d *Developer) GetPaid() {
	fmt.Printf("%s is getting paid\n", d.Name)
}

// These methods aren't relevant for a Developer, but the interface forces implementation
func (d *Developer) ManageOthers() {
	// Not applicable, but forced to implement
	fmt.Printf("%s doesn't manage others (empty implementation)\n", d.Name)
}

func (d *Developer) GenerateReports() {
	// Not applicable, but forced to implement
	fmt.Printf("%s doesn't generate reports (empty implementation)\n", d.Name)
}

func (d *Developer) AttendMeetings() {
	fmt.Printf("%s is attending a team standup\n", d.Name)
}

func (d *Developer) WorkOvertime() {
	fmt.Printf("%s is fixing bugs overtime\n", d.Name)
}

// Intern implements Worker but most methods are not applicable
type Intern struct {
	Name string
}

func (i *Intern) DoWork() {
	fmt.Printf("%s is learning and helping with simple tasks\n", i.Name)
}

func (i *Intern) TakeBreak() {
	fmt.Printf("%s is taking a study break\n", i.Name)
}

func (i *Intern) GetPaid() {
	fmt.Printf("%s is getting a small stipend\n", i.Name)
}

// Most of these methods make no sense for an Intern but they're forced to implement them
func (i *Intern) ManageOthers() {
	// Completely inappropriate, but forced to implement
	fmt.Printf("%s has no management responsibilities (empty implementation)\n", i.Name)
}

func (i *Intern) GenerateReports() {
	// Completely inappropriate, but forced to implement
	fmt.Printf("%s doesn't generate reports (empty implementation)\n", i.Name)
}

func (i *Intern) AttendMeetings() {
	fmt.Printf("%s is observing a meeting\n", i.Name)
}

func (i *Intern) WorkOvertime() {
	// Might not even be legal depending on jurisdiction
	fmt.Printf("%s is not supposed to work overtime (empty implementation)\n", i.Name)
}

// This function demonstrates Worker interface usage before ISP
func demonstrateWorkerBeforeISP() {
	workers := []Worker{
		&Manager{Name: "Alice"},
		&Developer{Name: "Bob"},
		&Intern{Name: "Charlie"},
	}

	fmt.Println("Morning routine:")
	for _, worker := range workers {
		worker.DoWork()
	}

	fmt.Println("\nMidday routine:")
	for _, worker := range workers {
		worker.TakeBreak()
	}

	fmt.Println("\nEnd of day routine:")
	for _, worker := range workers {
		worker.GetPaid()
	}

	fmt.Println("\nManagement activities:")
	for _, worker := range workers {
		// This is problematic since not all workers should manage others
		worker.ManageOthers()
		worker.GenerateReports()
	}

	fmt.Println("\nThe problem with this approach is:")
	fmt.Println("1. Developer and Intern must implement methods that don't apply to them")
	fmt.Println("2. This leads to empty implementations or inappropriate behavior")
	fmt.Println("3. Clients can call methods on objects that shouldn't have those methods")
	fmt.Println("4. If we add a new method to Worker, all implementations must be updated")
}
