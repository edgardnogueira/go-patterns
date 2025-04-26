package main

import (
	"fmt"
)

// BasicWorker defines the core behaviors all workers share
// This follows ISP by having only the essential methods
type BasicWorker interface {
	DoWork()
	TakeBreak()
	GetPaid()
}

// TeamMember defines behaviors for workers who participate in team activities
type TeamMember interface {
	AttendMeetings()
}

// OvertimeWorker defines behaviors for workers who can work extra hours
type OvertimeWorker interface {
	WorkOvertime()
}

// Manager defines behaviors specific to management roles
type ManagerRole interface {
	ManageOthers()
	GenerateReports()
}

// FullManager implements all interfaces, appropriate for a management role
type FullManager struct {
	Name string
}

func (m *FullManager) DoWork() {
	fmt.Printf("%s is delegating tasks\n", m.Name)
}

func (m *FullManager) TakeBreak() {
	fmt.Printf("%s is taking a coffee break\n", m.Name)
}

func (m *FullManager) GetPaid() {
	fmt.Printf("%s is getting paid a large salary\n", m.Name)
}

func (m *FullManager) ManageOthers() {
	fmt.Printf("%s is managing the team\n", m.Name)
}

func (m *FullManager) GenerateReports() {
	fmt.Printf("%s is generating performance reports\n", m.Name)
}

func (m *FullManager) AttendMeetings() {
	fmt.Printf("%s is attending a management meeting\n", m.Name)
}

func (m *FullManager) WorkOvertime() {
	fmt.Printf("%s is planning tomorrow's work\n", m.Name)
}

// FullDeveloper implements only the interfaces relevant for a developer
type FullDeveloper struct {
	Name string
}

func (d *FullDeveloper) DoWork() {
	fmt.Printf("%s is writing code\n", d.Name)
}

func (d *FullDeveloper) TakeBreak() {
	fmt.Printf("%s is browsing tech news\n", d.Name)
}

func (d *FullDeveloper) GetPaid() {
	fmt.Printf("%s is getting paid\n", d.Name)
}

func (d *FullDeveloper) AttendMeetings() {
	fmt.Printf("%s is attending a team standup\n", d.Name)
}

func (d *FullDeveloper) WorkOvertime() {
	fmt.Printf("%s is fixing bugs overtime\n", d.Name)
}

// Notice: No ManageOthers or GenerateReports methods as those don't apply

// FullIntern implements only the interfaces relevant for an intern
type FullIntern struct {
	Name string
}

func (i *FullIntern) DoWork() {
	fmt.Printf("%s is learning and helping with simple tasks\n", i.Name)
}

func (i *FullIntern) TakeBreak() {
	fmt.Printf("%s is taking a study break\n", i.Name)
}

func (i *FullIntern) GetPaid() {
	fmt.Printf("%s is getting a small stipend\n", i.Name)
}

func (i *FullIntern) AttendMeetings() {
	fmt.Printf("%s is observing a meeting\n", i.Name)
}

// Notice: No ManageOthers, GenerateReports or WorkOvertime as those don't apply

// MorningRoutine only requires BasicWorker interface
func MorningRoutine(workers []BasicWorker) {
	fmt.Println("Morning routine:")
	for _, worker := range workers {
		worker.DoWork()
	}
}

// MiddayRoutine only requires BasicWorker interface
func MiddayRoutine(workers []BasicWorker) {
	fmt.Println("\nMidday routine:")
	for _, worker := range workers {
		worker.TakeBreak()
	}
}

// EndOfDayRoutine only requires BasicWorker interface
func EndOfDayRoutine(workers []BasicWorker) {
	fmt.Println("\nEnd of day routine:")
	for _, worker := range workers {
		worker.GetPaid()
	}
}

// TeamMeeting only requires TeamMember interface
func TeamMeeting(members []TeamMember) {
	fmt.Println("\nTeam meeting:")
	for _, member := range members {
		member.AttendMeetings()
	}
}

// OvertimeWork only requires OvertimeWorker interface
func OvertimeWork(workers []OvertimeWorker) {
	fmt.Println("\nOvertime work:")
	for _, worker := range workers {
		worker.WorkOvertime()
	}
}

// ManagementActivities only requires ManagerRole interface
func ManagementActivities(managers []ManagerRole) {
	fmt.Println("\nManagement activities:")
	for _, manager := range managers {
		manager.ManageOthers()
		manager.GenerateReports()
	}
}

// This function demonstrates segregated interfaces
func demonstrateWorkerAfterISP() {
	// Create the different types of workers
	manager := &FullManager{Name: "Alice"}
	developer := &FullDeveloper{Name: "Bob"}
	intern := &FullIntern{Name: "Charlie"}

	// All are BasicWorkers
	basicWorkers := []BasicWorker{manager, developer, intern}
	MorningRoutine(basicWorkers)
	MiddayRoutine(basicWorkers)
	EndOfDayRoutine(basicWorkers)

	// All can attend meetings
	teamMembers := []TeamMember{manager, developer, intern}
	TeamMeeting(teamMembers)

	// Only Manager and Developer can work overtime
	overtimeWorkers := []OvertimeWorker{manager, developer}
	OvertimeWork(overtimeWorkers)

	// Only Manager can do management activities
	managers := []ManagerRole{manager}
	ManagementActivities(managers)

	// We cannot add intern to management activities or overtime
	// The following would cause compile errors:
	// overtimeWorkers = append(overtimeWorkers, intern)
	// managers = append(managers, intern)

	fmt.Println("\nThe benefits of this approach are:")
	fmt.Println("1. Clients only depend on the interfaces they need")
	fmt.Println("2. Each worker only implements the methods relevant to their role")
	fmt.Println("3. Type safety ensures workers are only used for appropriate tasks")
	fmt.Println("4. Adding new worker types or worker capabilities is simpler")
	fmt.Println("5. Adding a method to one interface doesn't affect other interfaces")
}
