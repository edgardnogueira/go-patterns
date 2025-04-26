package main

import (
	"errors"
	"testing"
)

// MockPaymentMethod is a test implementation of the PaymentMethodProcessor interface
type MockPaymentMethod struct {
	ProcessedAmount float64
	ShouldFail      bool
	ErrorToReturn   error
}

// ProcessPayment implements the PaymentMethodProcessor interface
func (m *MockPaymentMethod) ProcessPayment(amount float64) error {
	if m.ShouldFail {
		return m.ErrorToReturn
	}
	m.ProcessedAmount = amount
	return nil
}

func TestPaymentProcessorWithOCP(t *testing.T) {
	// Create processor
	processor := &PaymentProcessor{}
	
	// Test successful payment
	t.Run("successful payment", func(t *testing.T) {
		mockMethod := &MockPaymentMethod{}
		amount := 100.50
		
		err := processor.ProcessPayment(mockMethod, amount)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if mockMethod.ProcessedAmount != amount {
			t.Errorf("Expected processed amount %f, got %f", amount, mockMethod.ProcessedAmount)
		}
	})
	
	// Test failed payment
	t.Run("failed payment", func(t *testing.T) {
		mockMethod := &MockPaymentMethod{
			ShouldFail:    true,
			ErrorToReturn: errors.New("payment failed"),
		}
		amount := 200.75
		
		err := processor.ProcessPayment(mockMethod, amount)
		
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		
		if mockMethod.ProcessedAmount != 0 {
			t.Errorf("Expected no processed amount, got %f", mockMethod.ProcessedAmount)
		}
	})
	
	// Test adding a new payment method
	t.Run("adding new payment method", func(t *testing.T) {
		// Define a completely new payment method
		type GiftCardPayment struct {
			CardNumber  string
			Balance     float64
			processed   bool
		}
		
		// Implement PaymentMethodProcessor for the new method
		func(g *GiftCardPayment) ProcessPayment(amount float64) error {
			if g.Balance < amount {
				return errors.New("insufficient balance")
			}
			g.Balance -= amount
			g.processed = true
			return nil
		}
		
		// Create the new payment method
		giftCard := &GiftCardPayment{
			CardNumber: "GC12345",
			Balance:    50.0,
		}
		
		// Try to process a payment that exceeds the balance
		err := processor.ProcessPayment(giftCard, 75.0)
		if err == nil {
			t.Error("Expected insufficient balance error, got nil")
		}
		
		// Try to process a payment within the balance
		err = processor.ProcessPayment(giftCard, 30.0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if giftCard.Balance != 20.0 {
			t.Errorf("Expected balance 20.0, got %f", giftCard.Balance)
		}
		
		if !giftCard.processed {
			t.Error("Expected gift card to be processed")
		}
		
		// The OCP benefit: We added a new payment method without modifying
		// the PaymentProcessor or any existing code!
	})
}

func TestExistingPaymentMethods(t *testing.T) {
	// Create processor
	processor := &PaymentProcessor{}
	
	// Test CreditCardPayment
	t.Run("credit card payment", func(t *testing.T) {
		cc := &CreditCardPayment{
			CardNumber: "4111111111111111",
			CVV:        "123",
			Expiry:     "01/25",
		}
		
		err := processor.ProcessPayment(cc, 150.0)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test PayPalPayment
	t.Run("paypal payment", func(t *testing.T) {
		pp := &PayPalPayment{
			Email: "customer@example.com",
		}
		
		err := processor.ProcessPayment(pp, 75.50)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test CryptocurrencyPayment
	t.Run("cryptocurrency payment", func(t *testing.T) {
		crypto := &CryptocurrencyPayment{
			WalletAddress: "0x1234567890abcdef",
			Currency:      "Bitcoin",
		}
		
		err := processor.ProcessPayment(crypto, 0.5)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test invalid payment methods
	t.Run("invalid credit card", func(t *testing.T) {
		cc := &CreditCardPayment{
			// Missing required fields
		}
		
		err := processor.ProcessPayment(cc, 150.0)
		
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
	})
	
	t.Run("invalid paypal", func(t *testing.T) {
		pp := &PayPalPayment{
			// Missing email
		}
		
		err := processor.ProcessPayment(pp, 75.50)
		
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
	})
}
