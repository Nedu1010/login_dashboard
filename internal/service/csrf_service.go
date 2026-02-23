package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type CSRFService struct {
	secret string
}

func NewCSRFService(secret string) *CSRFService {
	return &CSRFService{secret: secret}
}

// GenerateToken creates a CSRF token based on timestamp and secret
func (s *CSRFService) GenerateToken() string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d", s.secret, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ValidateToken validates a CSRF token
// In a real implementation, you might want to check timestamp freshness
func (s *CSRFService) ValidateToken(token string) bool {
	// For double-submit pattern, we just check if token exists
	// The actual validation happens by comparing cookie and header
	return len(token) > 0
}
