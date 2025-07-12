package bcrypthashing

import "golang.org/x/crypto/bcrypt"

func BcryptHashing(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePasswordAndHash(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}