package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

var (
	certFile, pKeyFile *os.File
	certPath, pKeyPath string
)

func main1() {
	fmt.Println("Generating PEM files...")

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"H2P2F"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		panic(err)
	}

	var keyPEM bytes.Buffer
	if err := pem.Encode(&keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		panic(err)
	}

	certPath = "./crypto/cert.crypto"
	pKeyPath = "./crypto/key.crypto"
	certFile, err = os.Create(certPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := certFile.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := certFile.Write(certPEM.Bytes()); err != nil {
		panic(err)
	}

	pKeyFile, err = os.Create(pKeyPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := pKeyFile.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := pKeyFile.Write(keyPEM.Bytes()); err != nil {
		panic(err)
	}
	fmt.Println("cert.crypto and key.crypto generated successfully")
}
