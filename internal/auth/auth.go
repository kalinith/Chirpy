package auth
import(
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
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