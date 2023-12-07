package utils

import (
	"testing"
)

func Test_GenerateRandomString(t *testing.T) {
	length := 5
	randomStr := GenerateRandomString(length)

	if len(randomStr) != length {
		t.Fatalf("GenerateRandomString() => expected string length %d but got %d", length, len(randomStr))
	}
}

func Test_ValidateURL(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "invalid url",
			url:   "invalid-url",
			valid: false,
		},
		{
			name:  "valid url",
			url:   "https://github.com/sagar-jadhav",
			valid: true,
		},
		{
			name:  "invalid domain",
			url:   "https://" + LOCALHOST + ":3000/",
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateURL(test.url)
			if result != test.valid {
				t.Fatalf("ValidateURL() => expected %v but got %v", test.valid, result)
			}
		})
	}
}
