package strategy

import (
	"fmt"
	"strings"
)

// PaymentStrategy is the interface that all payment methods must implement
type PaymentStrategy interface {
	Pay(amount float64) string
	GetName() string
}

// CreditCardStrategy implements payment processing for credit cards
type CreditCardStrategy struct {
	Name     string
	CardNum  string
	CVV      string
	ExpMonth int
	ExpYear  int
}

// NewCreditCardStrategy creates a new credit card payment strategy
func NewCreditCardStrategy(name, cardNum, cvv string, expMonth, expYear int) *CreditCardStrategy {
	return &CreditCardStrategy{
		Name:     name,
		CardNum:  cardNum,
		CVV:      cvv,
		ExpMonth: expMonth,
		ExpYear:  expYear,
	}
}

// Pay processes a credit card payment
func (c *CreditCardStrategy) Pay(amount float64) string {
	// In a real implementation, this would integrate with a payment gateway
	return fmt.Sprintf("Paid %.2f using Credit Card (ending with %s)", 
		amount, 
		c.CardNum[len(c.CardNum)-4:])
}

// GetName returns the name of the payment method
func (c *CreditCardStrategy) GetName() string {
	return "Credit Card"
}

// PayPalStrategy implements payment processing for PayPal
type PayPalStrategy struct {
	Email    string
	Password string
}

// NewPayPalStrategy creates a new PayPal payment strategy
func NewPayPalStrategy(email, password string) *PayPalStrategy {
	return &PayPalStrategy{
		Email:    email,
		Password: password,
	}
}

// Pay processes a PayPal payment
func (p *PayPalStrategy) Pay(amount float64) string {
	// In a real implementation, this would integrate with PayPal's API
	return fmt.Sprintf("Paid %.2f using PayPal account: %s", 
		amount, 
		p.Email)
}

// GetName returns the name of the payment method
func (p *PayPalStrategy) GetName() string {
	return "PayPal"
}

// CryptoStrategy implements payment processing for cryptocurrency
type CryptoStrategy struct {
	CoinType  string
	WalletID  string
}

// NewCryptoStrategy creates a new cryptocurrency payment strategy
func NewCryptoStrategy(coinType, walletID string) *CryptoStrategy {
	return &CryptoStrategy{
		CoinType: coinType,
		WalletID: walletID,
	}
}

// Pay processes a cryptocurrency payment
func (c *CryptoStrategy) Pay(amount float64) string {
	// In a real implementation, this would integrate with a crypto payment processor
	return fmt.Sprintf("Paid %.2f using %s Wallet: %s", 
		amount, 
		c.CoinType,
		c.WalletID)
}

// GetName returns the name of the payment method
func (c *CryptoStrategy) GetName() string {
	return fmt.Sprintf("%s Cryptocurrency", c.CoinType)
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	Name     string
	Price    float64
	Quantity int
}

// ShoppingCart is the context that uses a payment strategy
type ShoppingCart struct {
	Items           []CartItem
	PaymentStrategy PaymentStrategy
}

// NewShoppingCart creates a new shopping cart
func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		Items: make([]CartItem, 0),
	}
}

// AddItem adds an item to the shopping cart
func (s *ShoppingCart) AddItem(item CartItem) {
	s.Items = append(s.Items, item)
}

// SetPaymentStrategy sets the payment strategy to use
func (s *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	s.PaymentStrategy = strategy
}

// GetTotal calculates the total price of all items in the cart
func (s *ShoppingCart) GetTotal() float64 {
	total := 0.0
	for _, item := range s.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Checkout processes payment for all items in the cart
func (s *ShoppingCart) Checkout() string {
	if s.PaymentStrategy == nil {
		return "Error: No payment method selected"
	}

	if len(s.Items) == 0 {
		return "Cart is empty"
	}

	amount := s.GetTotal()
	result := s.PaymentStrategy.Pay(amount)
	return result
}

// GetReceiptText generates a formatted receipt of the cart contents
func (s *ShoppingCart) GetReceiptText() string {
	if len(s.Items) == 0 {
		return "Cart is empty"
	}

	lines := []string{"Shopping Cart:"}
	lines = append(lines, "---------------------")
	for _, item := range s.Items {
		lines = append(lines, 
			fmt.Sprintf("%s (x%d) - $%.2f", item.Name, item.Quantity, item.Price * float64(item.Quantity)))
	}
	lines = append(lines, "---------------------")
	lines = append(lines, fmt.Sprintf("Total: $%.2f", s.GetTotal()))

	if s.PaymentStrategy != nil {
		lines = append(lines, fmt.Sprintf("Payment Method: %s", s.PaymentStrategy.GetName()))
	}
	
	return strings.Join(lines, "\n")
}
