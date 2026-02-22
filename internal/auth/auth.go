package auth

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"
)

// Dev credentials for sandbox/development
const DevPhone = "+15550000000"
const DevOTP = "123456"

// IsDev returns true when ENV is not production
func IsDev() bool {
	return os.Getenv("ENV") != "production" && os.Getenv("ENV") != "prod"
}

// Mock OTP for dev: "123456" works for any +1 number (legacy)
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
