package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUserManagerCreateUser(t *testing.T) {
	// Setup: Create a temporary test file
	tempFile, err := ioutil.TempFile("", "users-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName) // Clean up after test

	t.Run("user manager without SRP is harder to test", func(t *testing.T) {
		// This test demonstrates why code without SRP is harder to test:
		// 1. We have no way to mock the file operations
		// 2. We have no way to mock the email sending
		// 3. We have no way to verify that logging happened
		// 4. We can only test the end-to-end behavior

		// Create user manager
		userManager := NewUserManager()

		// Create a valid user
		validUser := User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Test creating a user
		// Note how we have to use real file operations here
		err := userManager.CreateUser(validUser)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// The only way to verify the operation is to check if the file exists
		// and contains the expected data
		_, err = os.Stat("users.json")
		if os.IsNotExist(err) {
			t.Error("users.json file should exist after CreateUser operation")
		}

		// Clean up the file
		os.Remove("users.json")

		// Limitations of this approach:
		// 1. Test is dependent on file system access
		// 2. Cannot verify email sending without a real email service
		// 3. Cannot isolate components for unit testing
		// 4. Cannot easily test error conditions in isolation
		
		t.Log("Note: This test only confirms basic functionality but can't properly test components in isolation")
		t.Log("The code without SRP forces us to test everything together and use real file operations")
	})

	t.Run("validation failure case", func(t *testing.T) {
		userManager := NewUserManager()

		// Create an invalid user (missing name)
		invalidUser := User{
			ID:       2,
			Name:     "", // Name is required
			Email:    "test@example.com",
			Password: "password123",
		}

		// Test validation
		err := userManager.CreateUser(invalidUser)
		if err == nil {
			t.Error("Expected validation error, got nil")
		}

		// We can test validation, but we have no way to confirm that:
		// 1. No file operations were attempted
		// 2. No email sending was attempted 
		// 3. Error was properly logged
		
		t.Log("Note: We can only test the external behavior, not internal interactions")
	})
}

func TestUserManagerTestability(t *testing.T) {
	t.Run("highlight testability issues without SRP", func(t *testing.T) {
		// This "test" only highlights the limitations of testing code without SRP

		// ISSUE 1: Can't mock file operations
		t.Log("ISSUE 1: Can't easily mock file operations")
		t.Log("- SaveUsersToFile directly writes to disk")
		t.Log("- LoadUsersFromFile directly reads from disk")
		t.Log("- Makes testing without side effects difficult")

		// ISSUE 2: Can't mock email sending
		t.Log("ISSUE 2: Can't easily mock email sending")
		t.Log("- SendWelcomeEmail is a method on UserManager")
		t.Log("- No way to inject a test double")
		t.Log("- Can't verify email was sent without actual sending")

		// ISSUE 3: Can't capture or verify logging
		t.Log("ISSUE 3: Can't capture or verify logging")
		t.Log("- Logging is embedded in UserManager")
		t.Log("- No way to check if correct messages were logged")

		// ISSUE 4: Testing error paths requires creating actual error conditions
		t.Log("ISSUE 4: Testing error paths requires creating actual error conditions")
		t.Log("- To test file write errors, need to make file system read-only")
		t.Log("- To test file read errors, need to corrupt files")

		// This is not a real test, just documentation
		t.Skip("This is not a real test, just documentation of testability issues")
	})
}
