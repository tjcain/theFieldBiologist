package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes stores the number of remember token bytes, it prevents
// using a number that is too small.
const RememberTokenBytes = 32

// RememberToken is a helper functionthat generates remember tokens of a size
// dictated by RememberTokenBytes
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

// Bytes will generate n random bytes or return an error
// it uses the crypto/rand package, so usage for tokens is safe.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slice of size nBytes and return a string
// that is the base64 URL encoded version of that byte slice.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
