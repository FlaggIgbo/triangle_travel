package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"triangle_travel/internal/auth"
)

// SendOTPRequest for POST /api/auth/send-otp
type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SendOTP creates/stores OTP for phone
func (h *Handlers) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone required"})
		return
	}
	phone := strings.TrimSpace(req.Phone)
	if !strings.HasPrefix(phone, "+1") || len(phone) != 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "US phone number required (+1XXXXXXXXXX)"})
		return
	}

	// Dev: accept dev phone without storing OTP
	if auth.IsDev() && phone == auth.DevPhone {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	// Store OTP (mock in dev: 123456 for any; prod: use Twilio)
	code := auth.GenerateOTP()
	expires := auth.OTPExpiry()
	_, err := h.DB.Exec(`
		INSERT INTO otp_codes (phone, code, expires_at) VALUES (?, ?, ?)
		ON CONFLICT(phone) DO UPDATE SET code = ?, expires_at = ?
	`, phone, code, expires, code, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	// TODO: In production, send via Twilio: h.twilio.Send(phone, "Your code is: "+code)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// VerifyOTPRequest for POST /api/auth/verify-otp
type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// VerifyOTP verifies code and returns JWT token
func (h *Handlers) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone and code required"})
		return
	}
	phone := strings.TrimSpace(req.Phone)
	code := strings.TrimSpace(req.Code)

	// Dev: accept dev phone + dev OTP (user is seeded on server start)
	if auth.IsDev() && phone == auth.DevPhone && code == auth.DevOTP {
		var userID int64
		err := h.DB.QueryRow("SELECT id FROM users WHERE phone = ?", phone).Scan(&userID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Dev user not seeded. Restart the server."})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		token, err := auth.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		expiresAt := auth.TokenExpiry()
		_, err = h.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	var storedCode string
	var expires interface{}
	err := h.DB.QueryRow("SELECT code, expires_at FROM otp_codes WHERE phone = ?", phone).Scan(&storedCode, &expires)
	if err == sql.ErrNoRows || code != storedCode {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
		return
	}

	// Get or create user
	var userID int64
	err = h.DB.QueryRow("SELECT id FROM users WHERE phone = ?", phone).Scan(&userID)
	if err == sql.ErrNoRows {
		res, err := h.DB.Exec("INSERT INTO users (phone) VALUES (?)", phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		userID, _ = res.LastInsertId()
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	expiresAt := auth.TokenExpiry()
	_, err = h.DB.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Delete used OTP
	h.DB.Exec("DELETE FROM otp_codes WHERE phone = ?", phone)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware extracts Bearer token and sets user_id in context
func (h *Handlers) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
		c.Abort()
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	var userID int64
	err := h.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ? AND expires_at > datetime('now')", token).Scan(&userID)
	if err == sql.ErrNoRows || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}
	c.Set("user_id", userID)
	c.Next()
}
