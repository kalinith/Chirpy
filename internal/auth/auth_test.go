package auth
import (
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestHashing(t *testing.T) {
	cases := []struct {
		input	string
	}{
		{
			input:    "hello world",
		},
                {
                        input:    "This is a test",
                },
                {
                        input:    "the Quick brown Fox",
                },
                {
                        input:    "jumped over the lazy dog",
                },
                {
                        input:    "Uno Duo Tres Quatro Sinco sinco ses",
                },
                {
                        input:    "thisisareallylongsentancewithnospaces",
                },
                {
                        input:    "oh",
                },
		// add more cases here
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		// get a hash for the password
		if err != nil {
			t.Errorf("failed to hash password '%s': %s",c.input, err)
		}
		// test if password can be compared to hash
		passmatch := CheckPasswordHash(hash, c.input)
		if passmatch != nil {
			t.Errorf("password(%s) doesn't match Hash(%s)", c.input, hash)
		}
	}

}

func TestJWTCreationAndValidation(t *testing.T) {
    // Create a random UUID
    userID := uuid.New()
    tokenSecret := "test-secret"
    
    // Test case 1: Valid token
    token, err := MakeJWT(userID, tokenSecret, time.Hour)
    if err != nil {
        t.Fatalf("Failed to create token: %v", err)
    }
    
    // Validate the token
    extractedID, err := ValidateJWT(token, tokenSecret)
    if err != nil {
        t.Fatalf("Failed to validate token: %v", err)
    }
    
    if extractedID != userID {
        t.Errorf("User ID mismatch. Got %v, want %v", extractedID, userID)
    }
    
    // Test case 2: Expired token
    expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour) // Expired 1 hour ago
    if err != nil {
        t.Fatalf("Failed to create expired token: %v", err)
    }
    
    _, err = ValidateJWT(expiredToken, tokenSecret)
    if err == nil {
        t.Error("Expected error for expired token, got nil")
    }
    
    // Test case 3: Wrong secret
    wrongSecret := "wrong-secret"
    _, err = ValidateJWT(token, wrongSecret)
    if err == nil {
        t.Error("Expected error for token with wrong secret, got nil")
    }
}