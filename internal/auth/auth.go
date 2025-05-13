package auth
import(
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	// Hash the password using the bcrypt.GenerateFromPassword function. Bcrypt
	// is a secure hash function that is intended for use with passwords.
	hashCost := 10 //do I need to move this to the cfg file at some point?
	if password == "" {
		return "", fmt.Errorf("invalid password supplied")
	}
	bytePassword := []byte(password)
	if len(bytePassword) > 72 {
		return "", fmt.Errorf("Password is too long")
	}
	hash, err := bcrypt.GenerateFromPassword(bytePassword, hashCost)
	if err != nil {
		log.Printf("unable to hash password '%s': %s", password, err)
		return "", fmt.Errorf("Unable to hash hash password")
	}
	return string(hash), nil
}

func CheckPasswordHash(hash, password string) error{
	// Use the bcrypt.CompareHashAndPassword function to compare the password
	// that the user entered in the HTTP request with the password that is stored in the database.
	if hash == "unset" {
		return fmt.Errorf("No password Set")
	}
	if password == "" {
		return fmt.Errorf("No password entered")
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	utcTime := time.Now().UTC()
	expiresIn, _ := time.ParseDuration("1h")
	expireTime := utcTime.Add(expiresIn)
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(utcTime),
		ExpiresAt: jwt.NewNumericDate(expireTime),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println(ss, err)
		return "", err
	}
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
	    // Make sure the signing method is what we expect
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	    }
	    
	    // Return the token secret for validation
	    return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// Check if the token is valid and get the claims
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        // Extract the user ID from the Subject field
        userID, err := uuid.Parse(claims.Subject)
        if err != nil {
            return uuid.Nil, err
        }
        return userID, nil
    }
    
    return uuid.Nil, fmt.Errorf("invalid token")
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("Auth Token missing")
	}
	bearer := strings.Split(authorization, " ")
	if len(bearer) > 2 {
		return "", fmt.Errorf("invalid Auth Token format")
	}
	return string(bearer[1]), nil
}

func MakeRefreshToken() (string, error) {
	bytestring := make([]byte, 32)
	rand.Read(bytestring)
	return hex.EncodeToString(bytestring), nil
}