package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatRequest for POST /api/chat
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// Chat handles AI chat (placeholder - integrate OpenAI/Anthropic later)
func (h *Handlers) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message required"})
		return
	}

	// Placeholder: echo back a simple response until AI is integrated
	reply := "Thanks for your question about \"" + req.Message + "\". The AI chat feature is coming soon. For now, you can use Triangle Travel to find stopover cities or add your flights in My Flights."
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
