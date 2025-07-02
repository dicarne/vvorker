package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"vvorker/conf"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// EncryptionConfig holds the configuration for encryption middleware
type EncryptionConfig struct {
	// Key is the encryption key (16, 24, or 32 bytes for AES-128, AES-192, or AES-256)
	Key []byte
	// HeaderName is the header that indicates if the request/response is encrypted
	HeaderName string
}

// DefaultEncryptionConfig returns a default encryption configuration
// Note: In production, the key should be stored securely and not hardcoded
func DefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		Key:        []byte(conf.AppConfigInstance.EncryptionKey), // Replace with your secure key
		HeaderName: "X-Encrypted-Data",
	}
}

// EncryptionMiddleware creates a new encryption middleware
func EncryptionMiddleware(config EncryptionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if not a POST, PUT, or PATCH request
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPut &&
			c.Request.Method != http.MethodGet &&
			c.Request.Method != http.MethodPatch {
			c.Next()
			return
		}

		// Check if request is encrypted
		if c.GetHeader(config.HeaderName) == "true" {
			// Decrypt request body
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
				return
			}

			if len(body) > 0 {
				decrypted, err := decrypt(body, config.Key)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to get request body"})
					return
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(decrypted))
			}

			// Process the request
			c.Next()
			return
		} else if conf.AppConfigInstance.EncryptionKey != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request is not right"})
		}

		c.Next()
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Encrypt encrypts plaintext using AES-GCM
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func decrypt(encrypted []byte, key []byte) ([]byte, error) {
	encrypted = encrypted[1 : len(encrypted)-1]
	logrus.Info("Decrypted data:", string(encrypted))
	ciphertext, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		logrus.Error("Failed to decode base64 data:", err)
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logrus.Error("Failed to create AES cipher:", err)
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logrus.Error("Failed to create AES-GCM:", err)
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logrus.Error("Failed to decrypt data:", err)
		return nil, err
	}

	return plaintext, nil
}

func CORSMiddlewaire(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With, X-Encrypted-Data")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Request.Header.Del("Origin")
		c.Next()
	}
}
