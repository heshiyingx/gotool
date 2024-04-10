package auth

import (
	"testing"
	"time"
)

// TestGenToken tests the GenToken function
func TestGenToken(t *testing.T) {
	// Define a payload for testing
	type testPayload struct {
		Name string
		Age  int
	}
	payload := testPayload{"Test", 30}

	// Call the GenToken function
	token, err := GenToken[testPayload](payload, "testSecret", 1)
	if err != nil {
		t.Errorf("GenToken failed with error: %v", err)
	}

	// Check if the token is not empty
	if token == "" {
		t.Errorf("Generated token is empty")
	}
}

// TestParseToken tests the ParseToken function
func TestParseToken(t *testing.T) {
	// Define a payload for testing
	type testPayload struct {
		Name string
		Age  int
	}
	payload := testPayload{"Test", 30}

	// Generate a token for testing
	token, _ := GenToken[testPayload](payload, "testSecret", time.Second*1)

	// Call the ParseToken function
	claims, err := ParseToken[testPayload](token, "testSecret")
	if err != nil {
		t.Errorf("ParseToken failed with error: %v", err)
	}

	// Check if the claims are not nil
	if claims == nil {
		t.Errorf("Parsed claims are nil")
	}

	// Check if the claims match the original payload
	if claims.JwtPayLoad != payload {
		t.Errorf("Parsed claims do not match the original payload")
	}
}
