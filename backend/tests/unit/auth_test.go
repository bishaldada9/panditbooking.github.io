package unit

import (
	"testing"

	"github.com/bees/hindu-ritual-platform/pkg/security"
)

func TestPasswordHashing(t *testing.T) {
	password := "SecurePass123!"
	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !security.CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password")
	}

	if security.CheckPassword(hash, "WrongPassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"Abc123!@", true},
		{"short1!", false},
		{"nouppercase123!", false},
		{"NOLOWERCASE123!", false},
		{"NoSpecialChar123", false},
		{"ValidPass123!", true},
	}

	for _, tt := range tests {
		err := security.ValidatePasswordStrength(tt.password)
		if tt.valid && err != nil {
			t.Errorf("Password %q should be valid but got error: %v", tt.password, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Password %q should be invalid but passed validation", tt.password)
		}
	}
}

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@test.co.uk", true},
		{"invalid-email", false},
		{"@missing-user.com", false},
		{"user@.com", false},
	}

	for _, tt := range tests {
		result := security.ValidateEmail(tt.email)
		if result != tt.valid {
			t.Errorf("Email %q validation = %v, want %v", tt.email, result, tt.valid)
		}
	}
}

func TestSecureTokenGeneration(t *testing.T) {
	token1, err := security.GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken failed: %v", err)
	}

	token2, err := security.GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken failed: %v", err)
	}

	if token1 == token2 {
		t.Error("Two generated tokens should be different")
	}

	if len(token1) < 32 {
		t.Errorf("Token length should be at least 32, got %d", len(token1))
	}
}

func TestInputSanitization(t *testing.T) {
	input := "<script>alert('xss')</script>"
	sanitized := security.SanitizeHTML(input)
	if sanitized == input {
		t.Error("SanitizeHTML should modify HTML tags")
	}
}

func TestPasswordValidation(t *testing.T) {
	password := "TestPass123!"
	if err := security.ValidatePasswordStrength(password); err != nil {
		t.Errorf("Strong password should pass validation: %v", err)
	}
}
