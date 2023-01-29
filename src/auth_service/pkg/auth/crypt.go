package auth

import "golang.org/x/crypto/bcrypt"

func HashString(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(hash, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
