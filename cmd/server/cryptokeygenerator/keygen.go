package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Generate RSA keys...")
	keys, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	fmt.Println("RSA keys generated successfully")

	privateKeyPath := "./crypto/private.rsa"
	publicKeyPath := "./crypto/public.rsa"

	fmt.Println("Saving RSA keys...")
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := privateKeyFile.Close(); err != nil {
			panic(err)
		}
	}()

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := publicKeyFile.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Writing RSA keys...")

	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keys),
	})
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&keys.PublicKey),
	})
	if _, err := privateKeyFile.Write(privPem); err != nil {
		panic(err)
	}
	if _, err := publicKeyFile.Write(pubPem); err != nil {
		panic(err)
	}
	fmt.Println("RSA keys saved successfully")

}
