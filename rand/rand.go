package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytesLen = 32

func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func RememberToken() (string, error) {
	return generateRandString(RememberTokenBytesLen)
}

func generateRandString(nBytes int) (string, error) {
	b, err := generateRandBytes(RememberTokenBytesLen)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateRandBytes(nBytes int) ([]byte, error) {
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
