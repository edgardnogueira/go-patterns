package valueobject

import (
	"errors"
	"fmt"
)

// Money is a value object that represents a monetary value with currency
type Money struct {
	amount   int64  // Amount in smallest currency unit (e.g., cents)
	currency string // Currency code (e.g., USD, EUR)
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// MustNewMoney creates a new Money value object and panics if an error occurs
func MustNewMoney(amount int64, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// Amount returns the amount
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency
func (m Money) Currency() string {
	return m.currency
}

// Add adds another Money value and returns a new Money object
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", 
			m.currency, other.currency)
	}
	
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts another Money value and returns a new Money object
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies: %s and %s", 
			m.currency, other.currency)
	}
	
	if m.amount < other.amount {
		return Money{}, errors.New("insufficient amount")
	}
	
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply multiplies the amount by a factor and returns a new Money object
func (m Money) Multiply(factor int) Money {
	return Money{
		amount:   m.amount * int64(factor),
		currency: m.currency,
	}
}

// Equals checks if two Money objects are equal
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns a string representation of the Money
func (m Money) String() string {
	// Format amount as dollars/euros/etc
	major := m.amount / 100
	minor := m.amount % 100
	
	return fmt.Sprintf("%s %d.%02d", m.currency, major, minor)
}
