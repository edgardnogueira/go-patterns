package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// PaymentProvider is an interface for processing payments
type PaymentProvider interface {
	// ProcessPayment processes a payment and returns a transaction ID
	ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResult, error)
	
	// RefundPayment refunds a payment by transaction ID
	RefundPayment(ctx context.Context, transactionID string, amount float64) error
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	OrderID       string
	CustomerID    string
	Amount        float64
	Currency      string
	PaymentMethod string
}

// PaymentResult represents the result of a payment processing
type PaymentResult struct {
	TransactionID string
	Success       bool
	ErrorMessage  string
	Timestamp     time.Time
}

// ExternalPaymentProvider is an implementation that connects to an external payment service
type ExternalPaymentProvider struct {
	apiKey     string
	apiURL     string
	timeoutSec int
}

// NewExternalPaymentProvider creates a new ExternalPaymentProvider
func NewExternalPaymentProvider(apiKey, apiURL string, timeoutSec int) PaymentProvider {
	return &ExternalPaymentProvider{
		apiKey:     apiKey,
		apiURL:     apiURL,
		timeoutSec: timeoutSec,
	}
}

// ProcessPayment processes a payment through the external payment service
func (p *ExternalPaymentProvider) ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResult, error) {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would make an HTTP request to the payment gateway
	
	log.Printf("Processing payment for order %s: %.2f %s via %s", 
		request.OrderID, request.Amount, request.Currency, request.PaymentMethod)
	
	// Simulate API call delay
	select {
	case <-time.After(time.Duration(100) * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	
	// Generate a fake transaction ID
	transactionID := fmt.Sprintf("txn_%s_%d", request.OrderID, time.Now().Unix())
	
	// In this example, we'll approve all payments
	return &PaymentResult{
		TransactionID: transactionID,
		Success:       true,
		Timestamp:     time.Now(),
	}, nil
}

// RefundPayment refunds a payment through the external payment service
func (p *ExternalPaymentProvider) RefundPayment(ctx context.Context, transactionID string, amount float64) error {
	// This is a simplified implementation for demonstration purposes
	// In a real implementation, we would make an HTTP request to the payment gateway
	
	log.Printf("Refunding payment %s: %.2f", transactionID, amount)
	
	// Simulate API call delay
	select {
	case <-time.After(time.Duration(100) * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return ctx.Err()
	}
	
	// In this example, we'll approve all refunds
	return nil
}

// MockPaymentProvider is a mock implementation for testing
type MockPaymentProvider struct {
	shouldFail     bool
	failureMessage string
}

// NewMockPaymentProvider creates a new MockPaymentProvider
func NewMockPaymentProvider(shouldFail bool, failureMessage string) PaymentProvider {
	return &MockPaymentProvider{
		shouldFail:     shouldFail,
		failureMessage: failureMessage,
	}
}

// ProcessPayment processes a payment using mock logic
func (p *MockPaymentProvider) ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResult, error) {
	if p.shouldFail {
		return &PaymentResult{
			Success:      false,
			ErrorMessage: p.failureMessage,
			Timestamp:    time.Now(),
		}, errors.New(p.failureMessage)
	}
	
	// Generate a fake transaction ID
	transactionID := fmt.Sprintf("mock_txn_%s_%d", request.OrderID, time.Now().Unix())
	
	return &PaymentResult{
		TransactionID: transactionID,
		Success:       true,
		Timestamp:     time.Now(),
	}, nil
}

// RefundPayment mocks a refund operation
func (p *MockPaymentProvider) RefundPayment(ctx context.Context, transactionID string, amount float64) error {
	if p.shouldFail {
		return errors.New(p.failureMessage)
	}
	
	return nil
}
