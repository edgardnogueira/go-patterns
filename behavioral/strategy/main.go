package strategy

import (
	"fmt"
)

// This file contains example usage of the strategy pattern

// ExampleStrategyPattern demonstrates the strategy pattern in action
func ExampleStrategyPattern() {
	// Create a shopping cart
	cart := NewShoppingCart()

	// Add items to the cart
	cart.AddItem(CartItem{Name: "Laptop", Price: 999.99, Quantity: 1})
	cart.AddItem(CartItem{Name: "Mouse", Price: 19.99, Quantity: 2})
	cart.AddItem(CartItem{Name: "Keyboard", Price: 59.99, Quantity: 1})

	// Show the cart contents
	fmt.Println("Initial cart:")
	fmt.Println(cart.GetReceiptText())
	fmt.Println()

	// Attempt to checkout without a payment method
	fmt.Println("Attempting checkout without payment method:")
	fmt.Println(cart.Checkout())
	fmt.Println()

	// Use credit card payment strategy
	fmt.Println("Using Credit Card payment strategy:")
	creditCard := NewCreditCardStrategy("John Smith", "1234567890123456", "123", 12, 2025)
	cart.SetPaymentStrategy(creditCard)
	fmt.Println(cart.GetReceiptText())
	fmt.Println(cart.Checkout())
	fmt.Println()

	// Switch to PayPal payment strategy
	fmt.Println("Switching to PayPal payment strategy:")
	paypal := NewPayPalStrategy("john.smith@example.com", "mypassword")
	cart.SetPaymentStrategy(paypal)
	fmt.Println(cart.GetReceiptText())
	fmt.Println(cart.Checkout())
	fmt.Println()

	// Switch to cryptocurrency payment strategy
	fmt.Println("Switching to Cryptocurrency payment strategy:")
	crypto := NewCryptoStrategy("Bitcoin", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	cart.SetPaymentStrategy(crypto)
	fmt.Println(cart.GetReceiptText())
	fmt.Println(cart.Checkout())
}
