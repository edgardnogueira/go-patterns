package decorator

// TextProcessor is the Component interface in the Decorator pattern.
// It defines the basic operations that can be performed on text.
type TextProcessor interface {
	// Process takes input text and returns processed text and an error if any.
	Process(text string) (string, error)
	
	// GetName returns the name of the processor.
	GetName() string
	
	// GetDescription returns information about what the processor does.
	GetDescription() string
}

// BasicTextProcessor is a ConcreteComponent in the Decorator pattern.
// It provides the core text processing functionality.
type BasicTextProcessor struct {
	name        string
	description string
}

// NewBasicTextProcessor creates a new BasicTextProcessor.
func NewBasicTextProcessor() *BasicTextProcessor {
	return &BasicTextProcessor{
		name:        "Basic Text Processor",
		description: "Performs basic text processing without any transformations",
	}
}

// Process implements the TextProcessor interface.
// For the basic processor, it simply returns the input text unchanged.
func (b *BasicTextProcessor) Process(text string) (string, error) {
	return text, nil
}

// GetName returns the name of the processor.
func (b *BasicTextProcessor) GetName() string {
	return b.name
}

// GetDescription returns information about what the processor does.
func (b *BasicTextProcessor) GetDescription() string {
	return b.description
}

// TextProcessorDecorator is the base Decorator in the Decorator pattern.
// It wraps a TextProcessor and delegates operations to it.
type TextProcessorDecorator struct {
	wrapped     TextProcessor
	name        string
	description string
}

// Process delegates the processing to the wrapped TextProcessor.
func (d *TextProcessorDecorator) Process(text string) (string, error) {
	return d.wrapped.Process(text)
}

// GetName returns the name of the decorator.
func (d *TextProcessorDecorator) GetName() string {
	return d.name
}

// GetDescription returns information about what the decorator does.
func (d *TextProcessorDecorator) GetDescription() string {
	return d.description
}

// GetWrappedName returns the name of the wrapped processor.
func (d *TextProcessorDecorator) GetWrappedName() string {
	return d.wrapped.GetName()
}

// GetProcessingChain returns a string representing the chain of processors.
func (d *TextProcessorDecorator) GetProcessingChain() string {
	if decorator, ok := d.wrapped.(*TextProcessorDecorator); ok {
		return d.name + " → " + decorator.GetProcessingChain()
	}
	return d.name + " → " + d.wrapped.GetName()
}
