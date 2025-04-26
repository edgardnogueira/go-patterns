// Package interfaces demonstrates idiomatic Go interface implementation patterns
package interfaces

import (
	"fmt"
	"sort"
	"strings"
)

// Duck Typing in Go
// ----------------
// Duck typing is summed up by the saying: "If it walks like a duck and quacks like a duck, then it's a duck."
// In Go, this means types implicitly satisfy interfaces by implementing their methods,
// without explicitly declaring the relationship.

// Quacker is an interface for things that can quack
type Quacker interface {
	Quack() string
}

// Duck is a concrete type that naturally quacks
type Duck struct {
	Name string
}

// Quack makes the duck quack
func (d Duck) Quack() string {
	return fmt.Sprintf("%s: Quack! Quack!", d.Name)
}

// Person doesn't seem like a natural quacker, but can implement the behavior
type Person struct {
	Name string
}

// Quack makes the person imitate a duck
// This means Person implicitly satisfies the Quacker interface
func (p Person) Quack() string {
	return fmt.Sprintf("%s: (imitating) Quaaack!", p.Name)
}

// DuckTypingBasics demonstrates basic duck typing in Go
func DuckTypingBasics() {
	// Create instances of different types
	duck := Duck{Name: "Donald"}
	person := Person{Name: "John"}
	
	// Both can be treated as Quackers because they have the Quack() method
	quackers := []Quacker{duck, person}
	
	fmt.Println("Different types that can quack:")
	for _, q := range quackers {
		fmt.Println(q.Quack())
	}
}

// Real-world examples of duck typing
// --------------------------------

// Stringable is an interface for anything that can be converted to a string
type Stringable interface {
	ToString() string
}

// Various types implementing ToString

type Product struct {
	ID    int
	Name  string
	Price float64
}

func (p Product) ToString() string {
	return fmt.Sprintf("Product(ID: %d, Name: %s, Price: $%.2f)", p.ID, p.Name, p.Price)
}

type User struct {
	Name  string
	Email string
}

func (u User) ToString() string {
	return fmt.Sprintf("User(%s, %s)", u.Name, u.Email)
}

type Order struct {
	OrderID   string
	Products  []Product
	Total     float64
	Completed bool
}

func (o Order) ToString() string {
	status := "Pending"
	if o.Completed {
		status = "Completed"
	}
	return fmt.Sprintf("Order(ID: %s, Items: %d, Total: $%.2f, Status: %s)",
		o.OrderID, len(o.Products), o.Total, status)
}

// PrintItems prints any slice of items that implement Stringable
func PrintItems(items []Stringable) {
	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item.ToString())
	}
}

// Duck typing with standard library interfaces
// -----------------------------------------

// SortableString is just a string that knows how to sort itself
type SortableString string

// We can implement sort.Interface methods on any type
type SortableStringSlice []SortableString

func (s SortableStringSlice) Len() int { 
	return len(s) 
}

func (s SortableStringSlice) Less(i, j int) bool { 
	// Custom sorting: case-insensitive alphabetic order
	return strings.ToLower(string(s[i])) < strings.ToLower(string(s[j]))
}

func (s SortableStringSlice) Swap(i, j int) { 
	s[i], s[j] = s[j], s[i] 
}

// CustomSortExample demonstrates implementing sort.Interface
func CustomSortExample() {
	fruits := SortableStringSlice{"banana", "Apple", "cherry", "Durian"}
	fmt.Println("Before sorting:", fruits)
	
	// sort.Sort works on anything implementing sort.Interface
	sort.Sort(fruits)
	fmt.Println("After sorting:", fruits)
}

// Duck typing with structural interfaces
// -----------------------------------

// DisplayItem is an interface that only requires Properties() method
type DisplayItem interface {
	Properties() map[string]string
}

// WebPage implements DisplayItem
type WebPage struct {
	URL     string
	Title   string
	Content string
}

func (w WebPage) Properties() map[string]string {
	return map[string]string{
		"Type":  "Web Page",
		"URL":   w.URL,
		"Title": w.Title,
	}
}

// File implements DisplayItem
type File struct {
	Name     string
	Size     int
	Location string
}

func (f File) Properties() map[string]string {
	return map[string]string{
		"Type":     "File",
		"Name":     f.Name,
		"Size":     fmt.Sprintf("%d bytes", f.Size),
		"Location": f.Location,
	}
}

// DisplayProperties renders properties for any DisplayItem
func DisplayProperties(item DisplayItem) {
	props := item.Properties()
	fmt.Println("Item Properties:")
	for key, value := range props {
		fmt.Printf("  %s: %s\n", key, value)
	}
}

// Combining duck typing with struct embedding
// ----------------------------------------

// Logger provides logging functionality
type Logger struct {
	Prefix string
}

func (l Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Prefix, message)
}

// Service uses logging capability via embedding
type Service struct {
	Name string
	Logger // Embed Logger struct
}

// Loggable is an interface for anything that can log
type Loggable interface {
	Log(message string)
}

// UseLogger demonstrates using a logger without tight coupling
func UseLogger(logger Loggable, message string) {
	logger.Log(message)
}

// Advantages of duck typing
// -----------------------

// 1. Ex-post facto interfaces: You can create interfaces for existing types

// Moveable is an interface for things that can move
type Moveable interface {
	Move(dx, dy float64)
	Position() (x, y float64)
}

// Point satisfies Moveable
type Point2D struct {
	X, Y float64
}

func (p *Point2D) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

func (p Point2D) Position() (x, y float64) {
	return p.X, p.Y
}

// 2. Adaptability: Existing libraries can be used with new interfaces

// Wrapper to adapt existing types to new interfaces
type ReaderToStringable struct {
	Name   string
	reader Reader // Some reader from an existing library
}

func (r ReaderToStringable) ToString() string {
	return fmt.Sprintf("Reader(%s)", r.Name)
}

// StructuralVsNominalTyping demonstrates the difference between Go's
// structural typing vs nominal typing in other languages
func StructuralVsNominalTyping() {
	fmt.Println("\nStructural vs Nominal Typing:")
	fmt.Println("In Go (structural typing), a type satisfies an interface")
	fmt.Println("if it implements all the interface's methods - regardless of whether")
	fmt.Println("the relationship was declared explicitly.")
	
	fmt.Println("\nIn languages with nominal typing (e.g., Java), a class must")
	fmt.Println("explicitly declare that it implements an interface, even if")
	fmt.Println("it already has all the required methods.")
	
	fmt.Println("\nGo's approach enables greater flexibility and better separation")
	fmt.Println("between interface definitions and implementations.")
}

// DuckTypingDemo demonstrates duck typing in Go
func DuckTypingDemo() {
	fmt.Println("============================================")
	fmt.Println("Duck Typing in Go Demo")
	fmt.Println("============================================")
	
	// Basic example of duck typing
	DuckTypingBasics()
	
	// Real-world example
	fmt.Println("\nDuck typing with business objects:")
	
	// Create instances of different types
	product := Product{ID: 1, Name: "Laptop", Price: 999.99}
	user := User{Name: "Alice", Email: "alice@example.com"}
	order := Order{
		OrderID:   "ORD-12345",
		Products:  []Product{product},
		Total:     999.99,
		Completed: false,
	}
	
	// We can create a slice of Stringable with different types
	items := []Stringable{product, user, order}
	
	// Pass the slice to a function expecting Stringable items
	PrintItems(items)
	
	// Example with standard library interfaces
	fmt.Println("\nDuck typing with standard library interfaces (sort.Interface):")
	CustomSortExample()
	
	// Example with structural interfaces
	fmt.Println("\nDuck typing with structural interfaces:")
	webpage := WebPage{
		URL:     "https://example.com",
		Title:   "Example Website",
		Content: "Welcome to our site",
	}
	
	file := File{
		Name:     "document.pdf",
		Size:     1024 * 1024,
		Location: "/documents/",
	}
	
	DisplayProperties(webpage)
	DisplayProperties(file)
	
	// Example with struct embedding
	fmt.Println("\nDuck typing with struct embedding:")
	service := Service{
		Name:   "AuthService",
		Logger: Logger{Prefix: "AUTH"},
	}
	
	// We can use Service anywhere a Loggable is expected
	UseLogger(service, "Starting authentication service")
	UseLogger(&Logger{Prefix: "TEST"}, "Direct logger message")
	
	// Show differences between structural and nominal typing
	StructuralVsNominalTyping()
	
	fmt.Println("\nBenefits of Duck Typing in Go:")
	fmt.Println("1. Decoupling of interface definitions from implementations")
	fmt.Println("2. Ability to create interfaces for existing code")
	fmt.Println("3. More flexible library design")
	fmt.Println("4. Simplified testing with mocks and stubs")
	fmt.Println("5. Emphasis on behavior rather than type hierarchy")
	fmt.Println("============================================")
}
