package hashpass

import (
	"golang.org/x/crypto/bcrypt"
)

func HashSHA256(password []byte) []byte {
	// Сложность хеширования пароля
	cost := 10
	hash, _ := bcrypt.GenerateFromPassword(password, cost)
	return hash
}

func VerifyUserPass(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte("$2a$16$2CrpqhVaubaS.HAcsbgnxeWI9LHTKklgUnD8h/mfVs9W.yuGSvHLS"), []byte(password)); err != nil {
		return false
	}
	return true
}
