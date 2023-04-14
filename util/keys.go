package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func GenerateRSAPPr(bites int, targetPath string) error {
	pr, err := rsa.GenerateKey(rand.Reader, bites)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(pr)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	err = os.WriteFile(targetPath+"/private.pem", pem.EncodeToMemory(priBlock), 0644)
	if err != nil {
		return err
	}
	log.Println("private key generated successfully")
	publicKey := &pr.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	err = os.WriteFile(targetPath+"/public.pem", pem.EncodeToMemory(publicBlock), 0644)
	if err != nil {
		return err
	}
	log.Println("public key generated successfully")
	return nil
}
