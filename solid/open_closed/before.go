package main

import (
	"fmt"
)

// PaymentType enum to identify payment methods
type PaymentType int

const (
	CreditCard PaymentType = iota
	PayPal
	BankTransfer
)

// Payment represents payment information
type Payment struct {
	Amount      float64
	PaymentType PaymentType
	CreditCardNumber string
	CreditCardCVV    string
	CreditCardExpiry string
	PayPalEmail      string
	BankAccountName  string
	BankAccountNumber string
	BankRoutingNumber string
}

// PaymentProcessor processes payments
// This violates the Open/Closed Principle because adding a new payment method
// requires modifying the ProcessPayment method
type PaymentProcessor struct{}

// ProcessPayment processes a payment based on its type
// This method needs to be modified every time we add a new payment method
func (p *PaymentProcessor) ProcessPayment(payment Payment) error {
	// Depending on the payment type, use different processing logic
	switch payment.PaymentType {
	case CreditCard:
		// Process credit card payment
		return p.processCreditCardPayment(payment)
	case PayPal:
		// Process PayPal payment
		return p.processPayPalPayment(payment)
	case BankTransfer:
		// Process bank transfer
		return p.processBankTransferPayment(payment)
	default:
		return fmt.Errorf("unknown payment type")
	}
}

// processCreditCardPayment processes a credit card payment
func (p *PaymentProcessor) processCreditCardPayment(payment Payment) error {
	// In a real application, this would connect to a payment gateway
	fmt.Printf("Processing credit card payment of $%.2f with card number %s\n", 
		payment.Amount, payment.CreditCardNumber)
	// Validate credit card details
	if payment.CreditCardNumber == "" || payment.CreditCardCVV == "" || payment.CreditCardExpiry == "" {
		return fmt.Errorf("invalid credit card details")
	}
	return nil
}

// processPayPalPayment processes a PayPal payment
func (p *PaymentProcessor) processPayPalPayment(payment Payment) error {
	// In a real application, this would connect to PayPal's API
	fmt.Printf("Processing PayPal payment of $%.2f to account %s\n", 
		payment.Amount, payment.PayPalEmail)
	// Validate PayPal details
	if payment.PayPalEmail == "" {
		return fmt.Errorf("invalid PayPal account")
	}
	return nil
}

// processBankTransferPayment processes a bank transfer
func (p *PaymentProcessor) processBankTransferPayment(payment Payment) error {
	// In a real application, this would connect to a banking API
	fmt.Printf("Processing bank transfer of $%.2f to account %s (routing: %s)\n", 
		payment.Amount, payment.BankAccountNumber, payment.BankRoutingNumber)
	// Validate bank account details
	if payment.BankAccountName == "" || payment.BankAccountNumber == "" || payment.BankRoutingNumber == "" {
		return fmt.Errorf("invalid bank account details")
	}
	return nil
}

// This function demonstrates the payment processor before applying OCP
func demonstratePaymentProcessorBeforeOCP() {
	processor := &PaymentProcessor{}

	// Process a credit card payment
	ccPayment := Payment{
		Amount:           100.00,
		PaymentType:      CreditCard,
		CreditCardNumber: "4111111111111111",
		CreditCardCVV:    "123",
		CreditCardExpiry: "01/25",
	}
	if err := processor.ProcessPayment(ccPayment); err != nil {
		fmt.Println("Error processing credit card payment:", err)
	}

	// Process a PayPal payment
	ppPayment := Payment{
		Amount:      75.50,
		PaymentType: PayPal,
		PayPalEmail: "customer@example.com",
	}
	if err := processor.ProcessPayment(ppPayment); err != nil {
		fmt.Println("Error processing PayPal payment:", err)
	}

	// Process a bank transfer
	btPayment := Payment{
		Amount:            200.00,
		PaymentType:       BankTransfer,
		BankAccountName:   "John Doe",
		BankAccountNumber: "123456789",
		BankRoutingNumber: "987654321",
	}
	if err := processor.ProcessPayment(btPayment); err != nil {
		fmt.Println("Error processing bank transfer:", err)
	}

	// If we want to add a new payment method (e.g., Cryptocurrency),
	// we would need to:
	// 1. Add a new PaymentType constant
	// 2. Add a new case to the switch statement in ProcessPayment
	// 3. Add a new processing method
	// This violates the Open/Closed Principle because we're modifying existing code
}
