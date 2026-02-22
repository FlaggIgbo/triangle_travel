package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Mock OTP for dev: "123456" works for any +1 number
const MockOTP = "123456"

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GenerateOTP() string {
	return MockOTP // In production, use random 6 digits and send via Twilio
}

func TokenExpiry() time.Time {
	return time.Now().Add(30 * 24 * time.Hour) // 30 days
}

func OTPExpiry() time.Time {
	return time.Now().Add(10 * time.Minute)
}
