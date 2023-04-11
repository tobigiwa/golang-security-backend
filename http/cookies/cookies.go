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
	ErrCookieValueTooLong = errors.New("cookie value too long")
	ErrCookieValueInvalid = errors.New("invalid cookie value")
	ErrInvalidValue       = errors.New("invalid cookie value")
)

func EncodeCookieValue(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	if len(cookie.String()) > 4096 {
		return ErrCookieValueTooLong
	}
	http.SetCookie(w, &cookie)
	return nil
}
//DecodeCookieValue returns a string of the decode cookie value or error
func DecodeCookieValue(r *http.Request) (string, error) {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		return "", http.ErrNoCookie
	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrCookieValueInvalid
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
	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encryptValue := aseGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	cookie.Value = string(encryptValue)

	return EncodeCookieValue(w, cookie)
}
func ReadEncryptedCookie(r *http.Request, secretKey []byte) (string, error) {
	encryptedValue, err := DecodeCookieValue(r)
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
	if expectedName != "cokkie" {
		return "", ErrInvalidValue
	}
	return value, nil
}
