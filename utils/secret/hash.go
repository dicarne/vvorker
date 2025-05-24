package secret

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

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

	fmt.Println("====================================")
	fmt.Printf("password: %s, hash: %s\n", password, hash)

	start := time.Now()
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	elapsed := time.Since(start)
	fmt.Printf("bcrypt.CompareHashAndPassword took %s\n", elapsed)

	result := err == nil
	cache[cacheKey] = result
	return result
}

func MD5(content string) string {
	data := []byte(content)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
