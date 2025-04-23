package main

import (
	"fmt"
	"github.com/edgardnogueira/go-patterns/behavioral/strategy"
)

func main() {
	fmt.Println("Strategy Pattern Example")
	fmt.Println("=========================")
	fmt.Println("This example demonstrates a shopping cart that can use different payment"
		+ "\nstrategies without knowing the details of how each payment method works.")
	fmt.Println()

	// Create a shopping cart
	cart := strategy.NewShoppingCart()

	// Add items to the cart
	cart.AddItem(strategy.CartItem{Name: "Laptop", Price: 999.99, Quantity: 1})
	cart.AddItem(strategy.CartItem{Name: "Mouse", Price: 19.99, Quantity: 2})
	cart.AddItem(strategy.CartItem{Name: "Keyboard", Price: 59.99, Quantity: 1})

	// Display the cart
	fmt.Println("\n🛒 Shopping Cart:")
	fmt.Println("----------------")
	fmt.Println(cart.GetReceiptText())

	// Attempt to checkout without a payment method
	fmt.Println("\n❌ Attempting checkout without a payment method:")
	fmt.Printf("  » %s\n", cart.Checkout())

	// Try different payment strategies
	fmt.Println("\n💳 Using Credit Card payment strategy:")
	creditCard := strategy.NewCreditCardStrategy("John Smith", "1234567890123456", "123", 12, 2025)
	cart.SetPaymentStrategy(creditCard)
	fmt.Println(cart.GetReceiptText())
	fmt.Printf("  » %s\n", cart.Checkout())

	fmt.Println("\n💰 Switching to PayPal payment strategy:")
	paypal := strategy.NewPayPalStrategy("john.smith@example.com", "mypassword")
	cart.SetPaymentStrategy(paypal)
	fmt.Println(cart.GetReceiptText())
	fmt.Printf("  » %s\n", cart.Checkout())

	fmt.Println("\n🪙 Switching to Cryptocurrency payment strategy:")
	crypto := strategy.NewCryptoStrategy("Bitcoin", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	cart.SetPaymentStrategy(crypto)
	fmt.Println(cart.GetReceiptText())
	fmt.Printf("  » %s\n", cart.Checkout())

	fmt.Println("\nThe Strategy pattern allows us to change the payment method at runtime")
	fmt.Println("without changing the shopping cart logic. Each payment strategy is")
	fmt.Println("encapsulated in its own class, making it easy to add new payment methods.")
}
