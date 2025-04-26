package entities

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		email     string
		password  string
		firstName string
		lastName  string
		wantErr   error
	}{
		{
			name:      "Valid user",
			username:  "testuser",
			email:     "test@example.com",
			password:  "password123",
			firstName: "Test",
			lastName:  "User",
			wantErr:   nil,
		},
		{
			name:      "Empty username",
			username:  "",
			email:     "test@example.com",
			password:  "password123",
			firstName: "Test",
			lastName:  "User",
			wantErr:   ErrEmptyUsername,
		},
		{
			name:      "Invalid email",
			username:  "testuser",
			email:     "invalid",
			password:  "password123",
			firstName: "Test",
			lastName:  "User",
			wantErr:   ErrInvalidEmailFormat,
		},
		{
			name:      "Password too short",
			username:  "testuser",
			email:     "test@example.com",
			password:  "short",
			firstName: "Test",
			lastName:  "User",
			wantErr:   ErrPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.username, tt.email, tt.password, tt.firstName, tt.lastName)
			
			// Check error
			if err != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// If we expect an error, no need to check the user
			if tt.wantErr != nil {
				if user != nil {
					t.Errorf("NewUser() returned non-nil user when error occurred")
				}
				return
			}
			
			// Check user fields
			if user.Username != tt.username {
				t.Errorf("User.Username = %v, want %v", user.Username, tt.username)
			}
			if user.Email != tt.email {
				t.Errorf("User.Email = %v, want %v", user.Email, tt.email)
			}
			if user.Password != tt.password {
				t.Errorf("User.Password = %v, want %v", user.Password, tt.password)
			}
			if user.FirstName != tt.firstName {
				t.Errorf("User.FirstName = %v, want %v", user.FirstName, tt.firstName)
			}
			if user.LastName != tt.lastName {
				t.Errorf("User.LastName = %v, want %v", user.LastName, tt.lastName)
			}
			if user.ID == "" {
				t.Error("User.ID is empty, expected UUID")
			}
			
			// Timestamps should be set
			if user.CreatedAt.IsZero() {
				t.Error("User.CreatedAt is zero, expected timestamp")
			}
			if user.UpdatedAt.IsZero() {
				t.Error("User.UpdatedAt is zero, expected timestamp")
			}
		})
	}
}

func TestUserUpdate(t *testing.T) {
	// Create a test user
	user, err := NewUser("originaluser", "original@example.com", "password123", "Original", "User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	
	// Store original timestamps for comparison
	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt
	
	// Wait a moment to ensure updated timestamp will be different
	time.Sleep(1 * time.Millisecond)
	
	// Update the user
	err = user.Update("newusername", "new@example.com", "New", "Name")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	
	// Check updated fields
	if user.Username != "newusername" {
		t.Errorf("User.Username = %v, want %v", user.Username, "newusername")
	}
	if user.Email != "new@example.com" {
		t.Errorf("User.Email = %v, want %v", user.Email, "new@example.com")
	}
	if user.FirstName != "New" {
		t.Errorf("User.FirstName = %v, want %v", user.FirstName, "New")
	}
	if user.LastName != "Name" {
		t.Errorf("User.LastName = %v, want %v", user.LastName, "Name")
	}
	
	// CreatedAt should remain the same
	if user.CreatedAt != originalCreatedAt {
		t.Errorf("User.CreatedAt changed after update, should remain constant")
	}
	
	// UpdatedAt should be updated
	if user.UpdatedAt == originalUpdatedAt {
		t.Errorf("User.UpdatedAt did not change after update")
	}
	
	// Test invalid email update
	err = user.Update("", "invalid", "", "")
	if err != ErrInvalidEmailFormat {
		t.Errorf("Update() with invalid email, error = %v, want %v", err, ErrInvalidEmailFormat)
	}
}

func TestUserUpdatePassword(t *testing.T) {
	// Create a test user
	user, err := NewUser("testuser", "test@example.com", "password123", "Test", "User")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	
	// Store original update timestamp for comparison
	originalUpdatedAt := user.UpdatedAt
	
	// Wait a moment to ensure updated timestamp will be different
	time.Sleep(1 * time.Millisecond)
	
	// Update password
	err = user.UpdatePassword("newpassword123")
	if err != nil {
		t.Fatalf("Failed to update password: %v", err)
	}
	
	// Check password was updated
	if user.Password != "newpassword123" {
		t.Errorf("User.Password = %v, want %v", user.Password, "newpassword123")
	}
	
	// UpdatedAt should be updated
	if user.UpdatedAt == originalUpdatedAt {
		t.Errorf("User.UpdatedAt did not change after password update")
	}
	
	// Test with password too short
	err = user.UpdatePassword("short")
	if err != ErrPasswordTooShort {
		t.Errorf("UpdatePassword() with short password, error = %v, want %v", err, ErrPasswordTooShort)
	}
}

func TestUserFullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		want      string
	}{
		{
			name:      "Both names provided",
			firstName: "John",
			lastName:  "Doe",
			want:      "John Doe",
		},
		{
			name:      "Only first name",
			firstName: "John",
			lastName:  "",
			want:      "John",
		},
		{
			name:      "Only last name",
			firstName: "",
			lastName:  "Doe",
			want:      "Doe",
		},
		{
			name:      "No names provided",
			firstName: "",
			lastName:  "",
			want:      "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create user with test names
			user, err := NewUser("testuser", "test@example.com", "password123", tt.firstName, tt.lastName)
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
			
			fullName := user.FullName()
			if fullName != tt.want {
				t.Errorf("FullName() = %v, want %v", fullName, tt.want)
			}
		})
	}
}
