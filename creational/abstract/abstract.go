package abstract

// Button is an abstract product
type Button interface {
	Render() string
	OnClick() string
}

// Checkbox is an abstract product
type Checkbox interface {
	Render() string
	Toggle() string
}

// GUIFactory is the abstract factory interface
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// ModernButton is a concrete product
type ModernButton struct{}

// Render returns a string representation of the modern button
func (b *ModernButton) Render() string {
	return "Rendered a modern style button"
}

// OnClick defines the button's click behavior
func (b *ModernButton) OnClick() string {
	return "Modern button clicked with smooth animation"
}

// VintageButton is a concrete product
type VintageButton struct{}

// Render returns a string representation of the vintage button
func (b *VintageButton) Render() string {
	return "Rendered a vintage style button"
}

// OnClick defines the button's click behavior
func (b *VintageButton) OnClick() string {
	return "Vintage button clicked with click sound"
}

// ModernCheckbox is a concrete product
type ModernCheckbox struct{}

// Render returns a string representation of the modern checkbox
func (c *ModernCheckbox) Render() string {
	return "Rendered a modern style checkbox"
}

// Toggle defines the checkbox's toggle behavior
func (c *ModernCheckbox) Toggle() string {
	return "Modern checkbox toggled with sliding animation"
}

// VintageCheckbox is a concrete product
type VintageCheckbox struct{}

// Render returns a string representation of the vintage checkbox
func (c *VintageCheckbox) Render() string {
	return "Rendered a vintage style checkbox"
}

// Toggle defines the checkbox's toggle behavior
func (c *VintageCheckbox) Toggle() string {
	return "Vintage checkbox toggled with mechanical sound"
}

// ModernGUIFactory is a concrete factory implementing GUIFactory
type ModernGUIFactory struct{}

// CreateButton creates a modern button
func (f *ModernGUIFactory) CreateButton() Button {
	return &ModernButton{}
}

// CreateCheckbox creates a modern checkbox
func (f *ModernGUIFactory) CreateCheckbox() Checkbox {
	return &ModernCheckbox{}
}

// VintageGUIFactory is a concrete factory implementing GUIFactory
type VintageGUIFactory struct{}

// CreateButton creates a vintage button
func (f *VintageGUIFactory) CreateButton() Button {
	return &VintageButton{}
}

// CreateCheckbox creates a vintage checkbox
func (f *VintageGUIFactory) CreateCheckbox() Checkbox {
	return &VintageCheckbox{}
}

// CreateGUIFactory returns a GUIFactory based on the style
func CreateGUIFactory(style string) GUIFactory {
	switch style {
	case "modern":
		return &ModernGUIFactory{}
	case "vintage":
		return &VintageGUIFactory{}
	default:
		return &ModernGUIFactory{} // Default to modern style
	}
}
