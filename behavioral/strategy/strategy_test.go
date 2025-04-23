package strategy

import (
	"strings"
	"testing"
)

func TestCreditCardStrategy(t *testing.T) {
	creditCard := NewCreditCardStrategy("John Smith", "1234567890123456", "123", 12, 2025)
	result := creditCard.Pay(100.50)

	expected := "Paid 100.50 using Credit Card (ending with 3456)"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	if creditCard.GetName() != "Credit Card" {
		t.Errorf("Expected name 'Credit Card', got '%s'", creditCard.GetName())
	}
}

func TestPayPalStrategy(t *testing.T) {
	paypal := NewPayPalStrategy("john.smith@example.com", "password")
	result := paypal.Pay(100.50)

	expected := "Paid 100.50 using PayPal account: john.smith@example.com"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	if paypal.GetName() != "PayPal" {
		t.Errorf("Expected name 'PayPal', got '%s'", paypal.GetName())
	}
}

func TestCryptoStrategy(t *testing.T) {
	crypto := NewCryptoStrategy("Bitcoin", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	result := crypto.Pay(0.005)

	// We just check if the result contains the right components
	if !strings.Contains(result, "Paid 0.005") || !strings.Contains(result, "Bitcoin Wallet") {
		t.Errorf("Expected result to contain 'Paid 0.005' and 'Bitcoin Wallet', got '%s'", result)
	}

	expectedName := "Bitcoin Cryptocurrency"
	if crypto.GetName() != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, crypto.GetName())
	}
}

func TestShoppingCart(t *testing.T) {
	// Create a cart and add items
	cart := NewShoppingCart()
	cart.AddItem(CartItem{Name: "Laptop", Price: 999.99, Quantity: 1})
	cart.AddItem(CartItem{Name: "Mouse", Price: 19.99, Quantity: 2})

	// Test cart total
	expectedTotal := 1039.97 // 999.99 + (19.99 * 2)
	total := cart.GetTotal()
	if total != expectedTotal {
		t.Errorf("Expected total %.2f, got %.2f", expectedTotal, total)
	}

	// Test checkout with no payment method
	result := cart.Checkout()
	if result != "Error: No payment method selected" {
		t.Errorf("Expected error message, got '%s'", result)
	}

	// Test with credit card
	creditCard := NewCreditCardStrategy("John Smith", "1234567890123456", "123", 12, 2025)
	cart.SetPaymentStrategy(creditCard)
	result = cart.Checkout()
	if !strings.Contains(result, "Paid 1039.97 using Credit Card") {
		t.Errorf("Expected result to contain 'Paid 1039.97 using Credit Card', got '%s'", result)
	}

	// Test receipt
	receipt := cart.GetReceiptText()
	if !strings.Contains(receipt, "Laptop (x1) - $999.99") {
		t.Errorf("Expected receipt to contain 'Laptop (x1) - $999.99', got '%s'", receipt)
	}
	if !strings.Contains(receipt, "Mouse (x2) - $39.98") {
		t.Errorf("Expected receipt to contain 'Mouse (x2) - $39.98', got '%s'", receipt)
	}
	if !strings.Contains(receipt, "Total: $1039.97") {
		t.Errorf("Expected receipt to contain 'Total: $1039.97', got '%s'", receipt)
	}
	if !strings.Contains(receipt, "Payment Method: Credit Card") {
		t.Errorf("Expected receipt to contain 'Payment Method: Credit Card', got '%s'", receipt)
	}

	// Test strategy swapping (change to PayPal)
	paypal := NewPayPalStrategy("john.smith@example.com", "password")
	cart.SetPaymentStrategy(paypal)
	result = cart.Checkout()
	if !strings.Contains(result, "Paid 1039.97 using PayPal") {
		t.Errorf("Expected result to contain 'Paid 1039.97 using PayPal', got '%s'", result)
	}
}