package utils

import (
	"encoding/base64"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

// LoadFile reads a file from the filesystem at the given path
// and returns its content as a byte slice.
func LoadFile(filepath string) ([]byte, error) {
	var data []byte

	log.Println("Opening file: ", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return data, err
	}

	defer file.Close()

	log.Println("Reading file")
	data, err = io.ReadAll(file)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Base64Encode encodes the given data to base64.
func Base64Encode(data string) string {
	log.Println("Encoding data to base64")
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode decodes the given base64 data.
func Base64Decode(data string) (string, error) {
	log.Println("Decoding data from base64")
	decoded, err := base64.StdEncoding.DecodeString(data)
	return string(decoded), err
}

// GenerateSessionId generates a new session id.
func GenerateSessionId() string {
	log.Println("Generating session id")
	return uuid.New().String()
}
