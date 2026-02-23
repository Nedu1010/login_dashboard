// Package crypto provides cryptographic utilities for password hashing and token generation.
//
// Security is critical! This package uses industry-standard algorithms:
// - bcrypt for password hashing (slow by design to resist brute-force)
// - crypto/rand for generating cryptographically secure random tokens
package crypto

import (
	"crypto/rand"     // Cryptographically secure random number generator
	"encoding/base64" // Encode binary data to text
	"fmt"             // Error formatting

	"golang.org/x/crypto/bcrypt" // bcrypt password hashing algorithm
)

// bcryptCost determines how many times the hashing algorithm is applied.
// Higher = more secure but slower. Cost of 12 is recommended (2^12 = 4096 rounds).
//
// Cost comparison:
// - 10 = ~100ms to hash (too fast, vulnerable to brute-force)
// - 12 = ~300ms to hash (good balance, recommended)
// - 14 = ~1200ms to hash (very secure but may slow down login)
const bcryptCost = 12

// HashPassword hashes a plain-text password using bcrypt.
//
// How bcrypt works:
// 1. Generates a random "salt" (random data mixed with password)
// 2. Applies hashing algorithm 2^cost times
// 3. Result includes salt + hash in one string
//
// Why bcrypt?
// - Slow by design (makes brute-force attacks expensive)
// - Each password gets a random salt (same password = different hash)
// - "Adaptive" - can increase cost over time as computers get faster
//
// Example:
//
//	plain: "myPassword123"
//	hash:  "$2a$12$8jF5Q.../XYZ..." (60 characters)
//
// The hash contains: algorithm version + cost + salt + actual hash
func HashPassword(password string) (string, error) {
	// GenerateFromPassword does the heavy lifting
	// Returns hash as []byte, so we convert to string
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		// %w wraps the error, preserving context
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain-text password.
//
// Returns:
// - nil if passwords match (success!)
// - error if passwords don't match or hash is invalid
//
// How it works:
// 1. Extracts the salt from the hashedPassword
// 2. Hashes the plaintext password with that same salt
// 3. Compares the two hashes
//
// This is used during login to verify the user entered the correct password.
//
// Example usage:
//
//	err := ComparePassword("$2a$12...", "myPassword123")
//	if err == nil {
//	    // Password is correct!
//	} else {
//	    // Password is wrong
//	}
//
// Security Note: This function runs in constant time to prevent timing attacks.
func ComparePassword(hashedPassword, password string) error {
	// bcrypt.CompareHashAndPassword does all the work
	// Returns nil if match, bcrypt.ErrMismatchedHashAndPassword if not
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateRandomToken generates a cryptographically secure random token.
//
// This is used for:
// - Refresh tokens (long-lived tokens stored in database)
// - Email verification tokens
// - Password reset tokens
//
// Parameters:
//   - length: Number of random bytes to generate (e.g., 32 bytes = 256 bits)
//
// Returns: Base64-encoded string (URL-safe, can be stored in database or sent in URLs)
//
// Why crypto/rand instead of math/rand?
// - crypto/rand uses OS's source of randomness (much more secure)
// - math/rand is predictable (uses a seed) - NEVER use it for security!
//
// Example:
//
//	token, _ := GenerateRandomToken(32)
//	// Result: "8jF5Q7Xm9pL2kN3oP6qR8sT4uV7wX1yZ..."
func GenerateRandomToken(length int) (string, error) {
	// Create a byte slice to hold the random data
	bytes := make([]byte, length)

	// Fill bytes with cryptographically secure random data
	// Read from crypto/rand (which uses /dev/urandom on Linux)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode to base64 for safe storage and transmission
	// URLEncoding is URL-safe (uses - and _ instead of + and /)
	return base64.URLEncoding.EncodeToString(bytes), nil
}
