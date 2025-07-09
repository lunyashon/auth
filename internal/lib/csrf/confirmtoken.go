package csrf

import (
	"bytes"
	"fmt"
	"math/rand"
)

func GenerateConfirmToken() string {
	var res bytes.Buffer

	for i := 0; i < 6; i++ {
		res.WriteString(fmt.Sprintf("%d", rand.Intn(10)))
	}

	return res.String()
}
