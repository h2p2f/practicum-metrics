// Package: cryptokeygenerator содержит в себе генератор RSA ключей.
// после генерации ключей, они сохраняются в папку crypto в формате PEM.
//
// Package: cryptokeygenerator contains an RSA key generator.
// after generating the keys, they are saved to the crypto folder in PEM format.
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
	// генерируем ключи
	// generate keys
	keys, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	fmt.Println("RSA keys generated successfully")
	// пути к файлам ключей
	// paths to key files
	privateKeyPath := "./crypto/private.rsa"
	publicKeyPath := "./crypto/public.rsa"

	fmt.Println("Saving RSA keys...")
	// создаем файлы ключей
	// create key files
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
	// записываем ключи в файлы
	// write keys to files
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
