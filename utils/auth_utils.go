package utils

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Decrypt(cypherText string) (string, error) {

	var data string

	privateKeyPath := os.Getenv(PRIVATE_KEY_PATH)
	log.Println("privateKeyPath: ", privateKeyPath)
	privateKey, err := LoadFile(privateKeyPath)
	if err != nil {
		return data, err
	}

	log.Println("Decoding private key")
	privateKeyPem, _ := pem.Decode(privateKey)
	if privateKey == nil {
		return data, err
	}

	key, err := x509.ParsePKCS8PrivateKey(privateKeyPem.Bytes)
	if err != nil {
		return data, err
	}

	pvtKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return data, fmt.Errorf("not a valid private key")
	}
	log.Println("Decrypting data")
	dataByte, err := rsa.DecryptOAEP(sha256.New(), nil, pvtKey, []byte(cypherText), nil)
	if err != nil {
		fmt.Println("--->", err)
		return data, err
	}

	log.Println("base64 Decoding data")
	data, err = Base64Decode(string(dataByte))
	if err != nil {
		return data, err
	}
	return data, nil
}

func HashPassword(pwd string) (string, error) {
	log.Println("Hashing password")
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hashPwd), err
}

func ComparePasswordHash(hash, pwd string) bool {
	log.Println("Comparing password hash")
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
