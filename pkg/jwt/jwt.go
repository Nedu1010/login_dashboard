// Package jwt provides utilities for creating and validating JWT (JSON Web Tokens).
//
// What is JWT?
// - A compact, URL-safe means of representing claims between two parties
// - Format: header.payload.signature (three base64-encoded parts separated by dots)
// - Example: eyJhbGc...  (header) . eyJzdWI... (payload) . SflKxw... (signature)
//
// Why use JWT?
// - Stateless: Server doesn't need to store session data
// - Self-contained: All info is in the token itself
// - Verifiable: Signature proves the token hasn't been tampered with
package jwt

import (
	"fmt"  // For error formatting
	"time" // For token expiration times

	"github.com/golang-jwt/jwt/v5" // Popular JWT library for Go
)

// Claims represents the data stored inside the JWT token.
// Think of this as the "payload" - the actual information the token carries.
//
// RegisteredClaims are standard JWT fields (ExpiresAt, IssuedAt, etc.)
// We embed it here so our Claims has all those fields automatically.
type Claims struct {
	UserID               int64  `json:"user_id"` // Custom claim: which user this token belongs to
	Email                string `json:"email"`   // Custom claim: user's email
	jwt.RegisteredClaims        // Embedded struct - adds ExpiresAt, IssuedAt, etc.
}

// Embedding Explained:
// jwt.RegisteredClaims has fields like ExpiresAt, IssuedAt, Issuer, etc.
// By embedding it (no field name), Claims automatically gets all those fields.
// You can access them like: claims.ExpiresAt

// GenerateAccessToken creates a short-lived JWT access token.
//
// Parameters:
//   - userID: The user's unique ID (stored in the token)
//   - email: The user's email (stored in the token)
//   - secret: Secret key used to sign the token (NEVER share this!)
//   - expiry: How long until the token expires (e.g., 5 minutes)
//
// Returns: A signed JWT string that can be verified later
//
// Security Note: The token is SIGNED, not ENCRYPTED.
// Anyone can read the payload, but only we can verify it's authentic.
// Never put sensitive data (passwords, credit cards) in JWT claims!
func GenerateAccessToken(userID int64, email string, secret string, expiry time.Duration) (string, error) {
	// Create the claims (payload data)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)), // Token expires after 'expiry' duration
			IssuedAt:  jwt.NewNumericDate(time.Now()),             // When was token created
			NotBefore: jwt.NewNumericDate(time.Now()),             // Token not valid before this time
		},
	}

	// Create a new token with our claims
	// SigningMethodHS256 = HMAC with SHA-256 (a symmetric signing algorithm)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key and return the string
	// The secret is converted to []byte because SignedString expects bytes
	// The resulting string is what we send to the client as a cookie
	return token.SignedString([]byte(secret))
}

// ValidateToken validates and parses a JWT token string.
//
// This function:
// 1. Checks if the token is properly formatted
// 2. Verifies the signature (proves it wasn't tampered with)
// 3. Checks if the token has expired
// 4. Returns the claims if everything is valid
//
// Parameters:
//   - tokenString: The JWT string from the client (from cookie or header)
//   - secret: The same secret key used to sign the token
//
// Returns: Parsed claims if valid, error otherwise
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	// Parse the token and extract claims
	// The third argument is a callback function that provides the secret key
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Security check: Verify the signing method is what we expect (HMAC)
		// This prevents attacks where someone changes the algorithm to "none"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return our secret key (converted to bytes) for verification
		return []byte(secret), nil
	})

	// Check if parsing failed
	if err != nil {
		// %w wraps the original error, preserving error chain for debugging
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Type assertion: Convert token.Claims to our *Claims type
	// ok will be false if the type assertion fails
	// token.Valid checks if the token hasn't expired and is properly signed
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	// If we reach here, something is wrong with the token
	return nil, fmt.Errorf("invalid token")
}
