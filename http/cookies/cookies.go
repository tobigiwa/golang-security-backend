package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrCookieValueTooLong  = errors.New("cookie value too long")
	ErrCookieValueInvalide = errors.New("invalid cookie value")
	ErrInvalidValue        = errors.New("invalid cookie value")
)

func EncodeCookieValue(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	if len(cookie.String()) > 4096 {
		return ErrCookieValueTooLong
	}
	http.SetCookie(w, &cookie)
	return nil
}

func DecodeCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrCookieValueInvalide
	}
	return string(value), nil
}

func WriteEncryptCookie(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}
	aseGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, aseGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encryptValue := aseGCM.Seal(nonce, nonce, []byte(text), nil)
	cookie.Value = string(encryptValue)

	return EncodeCookieValue(w, cookie)
}
func ReadEncryptedCookie(r *http.Request, name string, secretKey []byte) (string, error) {
	encryptedValue, err := DecodeCookieValue(r, name)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	aseGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aseGCM.NonceSize()
	if len(encryptedValue) < nonceSize {
		return "", err
	}

	nonce := encryptedValue[:nonceSize]
	value := encryptedValue[nonceSize:]

	plaintext, err := aseGCM.Open(nil, []byte(nonce), []byte(value), nil)
	if err != nil {
		return "", err
	}
	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidValue
	}
	if expectedName != name {
		return "", ErrInvalidValue
	}
	return value, nil
}
