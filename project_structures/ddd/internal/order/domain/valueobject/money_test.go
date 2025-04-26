package valueobject

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	// Test valid money creation
	money, err := NewMoney(1000, "USD")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if money.Amount() != 1000 {
		t.Errorf("Expected amount 1000, got %d", money.Amount())
	}
	
	if money.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", money.Currency())
	}
	
	// Test invalid money creation - negative amount
	_, err = NewMoney(-100, "USD")
	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
	
	// Test invalid money creation - empty currency
	_, err = NewMoney(100, "")
	if err == nil {
		t.Error("Expected error for empty currency, got nil")
	}
}

func TestMustNewMoney(t *testing.T) {
	// Test valid case (shouldn't panic)
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustNewMoney panicked unexpectedly: %v", r)
			}
		}()
		
		money := MustNewMoney(1000, "USD")
		
		if money.Amount() != 1000 {
			t.Errorf("Expected amount 1000, got %d", money.Amount())
		}
	}()
	
	// Test invalid case (should panic)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected MustNewMoney to panic with invalid args, but it didn't")
			}
		}()
		
		MustNewMoney(-100, "USD") // This should panic
	}()
}

func TestMoneyAdd(t *testing.T) {
	money1 := MustNewMoney(1000, "USD") // $10.00
	money2 := MustNewMoney(500, "USD")  // $5.00
	
	// Test adding same currency
	sum, err := money1.Add(money2)
	if err != nil {
		t.Fatalf("Failed to add money: %v", err)
	}
	
	if sum.Amount() != 1500 {
		t.Errorf("Expected sum 1500, got %d", sum.Amount())
	}
	
	if sum.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", sum.Currency())
	}
	
	// Test adding different currencies (should fail)
	money3 := MustNewMoney(800, "EUR") // €8.00
	_, err = money1.Add(money3)
	if err == nil {
		t.Error("Expected error when adding different currencies, got nil")
	}
}

func TestMoneySubtract(t *testing.T) {
	money1 := MustNewMoney(1000, "USD") // $10.00
	money2 := MustNewMoney(300, "USD")  // $3.00
	
	// Test subtracting same currency
	diff, err := money1.Subtract(money2)
	if err != nil {
		t.Fatalf("Failed to subtract money: %v", err)
	}
	
	if diff.Amount() != 700 {
		t.Errorf("Expected difference 700, got %d", diff.Amount())
	}
	
	if diff.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", diff.Currency())
	}
	
	// Test subtracting different currencies (should fail)
	money3 := MustNewMoney(200, "EUR") // €2.00
	_, err = money1.Subtract(money3)
	if err == nil {
		t.Error("Expected error when subtracting different currencies, got nil")
	}
	
	// Test insufficient amount (should fail)
	_, err = money2.Subtract(money1)
	if err == nil {
		t.Error("Expected error when subtracting larger amount, got nil")
	}
}

func TestMoneyMultiply(t *testing.T) {
	money := MustNewMoney(500, "USD") // $5.00
	
	// Test multiplying by positive factor
	result := money.Multiply(3)
	
	if result.Amount() != 1500 {
		t.Errorf("Expected amount 1500, got %d", result.Amount())
	}
	
	if result.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", result.Currency())
	}
	
	// Test multiplying by zero
	result = money.Multiply(0)
	
	if result.Amount() != 0 {
		t.Errorf("Expected amount 0, got %d", result.Amount())
	}
}

func TestMoneyEquals(t *testing.T) {
	money1 := MustNewMoney(1000, "USD")
	money2 := MustNewMoney(1000, "USD")
	money3 := MustNewMoney(1000, "EUR")
	money4 := MustNewMoney(1500, "USD")
	
	// Test equal money objects
	if !money1.Equals(money2) {
		t.Error("Expected money1 and money2 to be equal, but they weren't")
	}
	
	// Test different currencies
	if money1.Equals(money3) {
		t.Error("Expected money1 and money3 to be different, but they were equal")
	}
	
	// Test different amounts
	if money1.Equals(money4) {
		t.Error("Expected money1 and money4 to be different, but they were equal")
	}
}

func TestMoneyString(t *testing.T) {
	money1 := MustNewMoney(1050, "USD") // $10.50
	str1 := money1.String()
	if str1 != "USD 10.50" {
		t.Errorf("Expected string 'USD 10.50', got '%s'", str1)
	}
	
	money2 := MustNewMoney(5, "EUR") // €0.05
	str2 := money2.String()
	if str2 != "EUR 0.05" {
		t.Errorf("Expected string 'EUR 0.05', got '%s'", str2)
	}
	
	money3 := MustNewMoney(200, "JPY") // ¥2.00
	str3 := money3.String()
	if str3 != "JPY 2.00" {
		t.Errorf("Expected string 'JPY 2.00', got '%s'", str3)
	}
}
