package user_test

import (
	"testing"

	"github.com/P1coFly/VK_Authorization/internal/user"
)

func TestIsEmailValid(t *testing.T) {
	testCases := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"artem3m@mail.ru", true},
		{"test@t", true},
		{"test@", false},
	}

	for _, tc := range testCases {
		t.Run(tc.email, func(t *testing.T) {
			err := user.IsEmailValid(tc.email)
			if err != nil && tc.expected {
				t.Errorf("Expected %s to be a valid email, but got error: %v", tc.email, err)
			}
			if err == nil && !tc.expected {
				t.Errorf("Expected %s to be an invalid email, but got no error", tc.email)
			}
		})
	}
}

func TestPasswordCheckStatus(t *testing.T) {
	testCases := []struct {
		password string
		expected string
	}{
		{"12345", "weak"},
		{"12345123", "weak"},
		{"qwerqty", "weak"},
		{"qwerqty$", "good"},
		{"12345123$", "good"},
		{"password123", "good"},
		{"P@ssw0rd!", "perfect"},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			status := user.PasswordCheckStatus(tc.password)
			if status != tc.expected {
				t.Errorf("Expected %s to have '%s' status, but got %s", tc.password, tc.expected, status)
			}
		})
	}
}
