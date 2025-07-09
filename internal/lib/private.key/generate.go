/*
*
Update Private and Public key (6 month)
and write to .env file in this project
*/
package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	if err := WriteEnv(privateKeyPem, publicKeyPem); err != nil {
		panic(err)
	}
}

func WriteEnv(privateKey, publicKey []byte) error {

	var (
		path   string
		writer = make(map[string]string)
	)

	flag.StringVar(&path, "env", "", "path to env file")
	flag.Parse()

	if _, err := os.Stat(path); err != nil {
		return err
	}

	if file, err := os.Open(path); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			if params := strings.Split(line, "="); len(params) == 2 {
				key := strings.TrimSpace(params[0])
				val := strings.TrimSpace(params[1])
				writer[key] = val
			}
		}
		file.Close()
	} else {
		return err
	}

	writer["PRIVATE_KEY"] = strings.ReplaceAll(string(privateKey), "\n", "\\n")
	writer["PUBLIC_KEY"] = strings.ReplaceAll(string(publicKey), "\n", "\\n")

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	for key, value := range writer {
		if _, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
			return err
		}
	}

	return nil
}
