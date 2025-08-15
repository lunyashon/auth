package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token, pepper string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
