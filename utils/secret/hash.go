package secret

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

var cache = make(map[string]bool)

func CheckPasswordHash(password, hash string) bool {
	cacheKey := password + hash
	if result, ok := cache[cacheKey]; ok {
		return result
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	result := err == nil
	cache[cacheKey] = result
	return result
}

func MD5(content string) string {
	data := []byte(content)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
