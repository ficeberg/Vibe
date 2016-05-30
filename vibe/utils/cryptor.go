package utils

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
	"io"
	"strings"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

type SaltAuth struct{}

func (s *SaltAuth) Gen(password string) (string, string, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", "", err
	}

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return "", "", err
	}

	return strings.ToUpper(hex.EncodeToString(hash)), strings.ToUpper(hex.EncodeToString(salt)), nil
}

func (s *SaltAuth) Check(password, salt, hash string) bool {
	saltHex, _ := hex.DecodeString(salt)

	userhash, err := scrypt.Key([]byte(password), saltHex, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return false
	}

	return strings.ToUpper(hex.EncodeToString(userhash)) == hash
}
