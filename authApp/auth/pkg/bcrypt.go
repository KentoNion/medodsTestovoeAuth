package pkg

import "golang.org/x/crypto/bcrypt"

type Hash string

func HashToken(token Refresh) (Hash, error) { //функция перегоняющая токен в bcrypt hash
	tokenStr := string(token)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(tokenStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Hash(string(hashedBytes)), nil
}
