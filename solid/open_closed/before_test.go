package main

import (
	"testing"
)

func TestPaymentProcessorWithoutOCP(t *testing.T) {
	processor := &PaymentProcessor{}
	
	// Test credit card payment
	t.Run("credit card payment", func(t *testing.T) {
		payment := Payment{
			Amount:           100.00,
			PaymentType:      CreditCard,
			CreditCardNumber: "4111111111111111",
			CreditCardCVV:    "123",
			CreditCardExpiry: "01/25",
		}
		
		err := processor.ProcessPayment(payment)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test PayPal payment
	t.Run("paypal payment", func(t *testing.T) {
		payment := Payment{
			Amount:      75.50,
			PaymentType: PayPal,
			PayPalEmail: "customer@example.com",
		}
		
		err := processor.ProcessPayment(payment)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test bank transfer
	t.Run("bank transfer", func(t *testing.T) {
		payment := Payment{
			Amount:            200.00,
			PaymentType:       BankTransfer,
			BankAccountName:   "John Doe",
			BankAccountNumber: "123456789",
			BankRoutingNumber: "987654321",
		}
		
		err := processor.ProcessPayment(payment)
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	
	// Test invalid payment type
	t.Run("invalid payment type", func(t *testing.T) {
		payment := Payment{
			Amount:      150.00,
			PaymentType: 999, // Invalid type
		}
		
		err := processor.ProcessPayment(payment)
		
		if err == nil {
			t.Error("Expected error for invalid payment type, got nil")
		}
	})
	
	t.Run("cannot add cryptocurrency payment without modifying processor", func(t *testing.T) {
		// This test demonstrates the limitation of the approach without OCP
		
		// We can't easily add a Cryptocurrency payment type without:
		// 1. Adding a new PaymentType constant
		// 2. Modifying the ProcessPayment method to handle the new type
		// 3. Adding a new processing method
		
		t.Log("To add 'Cryptocurrency' as a payment method, we would need to:")
		t.Log("1. Add a new constant: Cryptocurrency PaymentType = iota + existing count")
		t.Log("2. Modify ProcessPayment to add a case for Cryptocurrency")
		t.Log("3. Create a new processCryptocurrencyPayment method")
		t.Log("4. Update the Payment struct to include cryptocurrency fields")
		
		// Since this is just a test highlighting limitations,
		// we'll just skip the actual implementation
		t.Skip("This test is just meant to demonstrate OCP violation")
	})
}

func TestPaymentProcessorExtensibilityIssues(t *testing.T) {
	t.Run("highlight issues without OCP", func(t *testing.T) {
		// This "test" highlights the limitations of code without OCP
		
		// ISSUE 1: Cannot add new payment methods without modifying existing code
		t.Log("ISSUE 1: Cannot add new payment methods without modifying existing code")
		t.Log("- The switch statement in ProcessPayment must be modified for each new payment type")
		t.Log("- This could introduce bugs in existing payment processing paths")
		
		// ISSUE 2: Testing becomes more complex as the switch grows
		t.Log("ISSUE 2: Testing complexity increases with each new payment type")
		t.Log("- Each time a new payment type is added, all existing tests must be rerun")
		t.Log("- Cannot easily isolate testing for each payment method")
		
		// ISSUE 3: PaymentType enum creates tight coupling
		t.Log("ISSUE 3: Tight coupling through the PaymentType enum")
		t.Log("- All payment methods must be known at compile time")
		t.Log("- Cannot dynamically register new payment methods")
		
		// ISSUE 4: Payment struct becomes bloated
		t.Log("ISSUE 4: Payment struct bloat")
		t.Log("- The Payment struct contains fields for all possible payment types")
		t.Log("- Most fields are irrelevant for any specific payment")
		
		// This is not a real test, just documentation
		t.Skip("This is not a real test, just documentation of extensibility issues")
	})
}
