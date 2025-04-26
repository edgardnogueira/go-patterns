package main

import (
	"fmt"
)

// PaymentMethodProcessor is an interface that all payment methods must implement
// This allows us to add new payment methods without modifying existing code
type PaymentMethodProcessor interface {
	ProcessPayment(amount float64) error
}

// PaymentProcessor processes payments using different payment methods
// This follows the Open/Closed Principle because it can process any payment method
// that implements the PaymentMethodProcessor interface without modification
type PaymentProcessor struct{}

// ProcessPayment processes a payment using the provided payment method
// This method doesn't need to change when we add new payment methods
func (p *PaymentProcessor) ProcessPayment(method PaymentMethodProcessor, amount float64) error {
	return method.ProcessPayment(amount)
}

// CreditCardPayment represents a credit card payment method
type CreditCardPayment struct {
	CardNumber string
	CVV        string
	Expiry     string
}

// ProcessPayment processes a credit card payment
func (c *CreditCardPayment) ProcessPayment(amount float64) error {
	// Validate credit card details
	if c.CardNumber == "" || c.CVV == "" || c.Expiry == "" {
		return fmt.Errorf("invalid credit card details")
	}
	// In a real application, this would connect to a payment gateway
	fmt.Printf("Processing credit card payment of $%.2f with card number %s\n", 
		amount, c.CardNumber)
	return nil
}

// PayPalPayment represents a PayPal payment method
type PayPalPayment struct {
	Email string
}

// ProcessPayment processes a PayPal payment
func (p *PayPalPayment) ProcessPayment(amount float64) error {
	// Validate PayPal account
	if p.Email == "" {
		return fmt.Errorf("invalid PayPal account")
	}
	// In a real application, this would connect to PayPal's API
	fmt.Printf("Processing PayPal payment of $%.2f to account %s\n", 
		amount, p.Email)
	return nil
}

// BankTransferPayment represents a bank transfer payment method
type BankTransferPayment struct {
	AccountName   string
	AccountNumber string
	RoutingNumber string
}

// ProcessPayment processes a bank transfer
func (b *BankTransferPayment) ProcessPayment(amount float64) error {
	// Validate bank account details
	if b.AccountName == "" || b.AccountNumber == "" || b.RoutingNumber == "" {
		return fmt.Errorf("invalid bank account details")
	}
	// In a real application, this would connect to a banking API
	fmt.Printf("Processing bank transfer of $%.2f to account %s (routing: %s)\n", 
		amount, b.AccountNumber, b.RoutingNumber)
	return nil
}

// CryptocurrencyPayment represents a cryptocurrency payment method
// This is a new payment method that can be added without modifying existing code
type CryptocurrencyPayment struct {
	WalletAddress string
	Currency      string
}

// ProcessPayment processes a cryptocurrency payment
func (c *CryptocurrencyPayment) ProcessPayment(amount float64) error {
	// Validate cryptocurrency details
	if c.WalletAddress == "" || c.Currency == "" {
		return fmt.Errorf("invalid cryptocurrency details")
	}
	// In a real application, this would connect to a cryptocurrency payment processor
	fmt.Printf("Processing %s cryptocurrency payment of $%.2f to wallet %s\n", 
		c.Currency, amount, c.WalletAddress)
	return nil
}

// This function demonstrates the payment processor after applying OCP
func demonstratePaymentProcessorAfterOCP() {
	processor := &PaymentProcessor{}

	// Process a credit card payment
	ccPayment := &CreditCardPayment{
		CardNumber: "4111111111111111",
		CVV:        "123",
		Expiry:     "01/25",
	}
	if err := processor.ProcessPayment(ccPayment, 100.00); err != nil {
		fmt.Println("Error processing credit card payment:", err)
	}

	// Process a PayPal payment
	ppPayment := &PayPalPayment{
		Email: "customer@example.com",
	}
	if err := processor.ProcessPayment(ppPayment, 75.50); err != nil {
		fmt.Println("Error processing PayPal payment:", err)
	}

	// Process a bank transfer
	btPayment := &BankTransferPayment{
		AccountName:   "John Doe",
		AccountNumber: "123456789",
		RoutingNumber: "987654321",
	}
	if err := processor.ProcessPayment(btPayment, 200.00); err != nil {
		fmt.Println("Error processing bank transfer:", err)
	}

	// Process a cryptocurrency payment - a new payment method
	// No changes to the PaymentProcessor are needed!
	cryptoPayment := &CryptocurrencyPayment{
		WalletAddress: "0x1234567890abcdef",
		Currency:      "Bitcoin",
	}
	if err := processor.ProcessPayment(cryptoPayment, 0.5); err != nil {
		fmt.Println("Error processing cryptocurrency payment:", err)
	}

	// To add another payment method (e.g., ApplePay, GooglePay),
	// we only need to:
	// 1. Create a new struct
	// 2. Implement the PaymentMethodProcessor interface
	// We DON'T need to modify the PaymentProcessor class or any existing code
	// This follows the Open/Closed Principle
}
