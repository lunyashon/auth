package csrf

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateResetToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
