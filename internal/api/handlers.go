package api

import (
	"net/http"
	"triangle_travel/internal/db"
	"triangle_travel/internal/flights"

	"github.com/gin-gonic/gin"
)

// Handlers holds dependencies
type Handlers struct {
	DB *db.DB
}

// SearchRequest for triangle travel
type SearchRequest struct {
	Start     string `json:"start" form:"start" binding:"required"`
	End       string `json:"end" form:"end" binding:"required"`
	StartDate string `json:"startDate" form:"startDate" binding:"required"`
	EndDate   string `json:"endDate" form:"endDate" binding:"required"`
	Cabin     string `json:"cabin" form:"cabin"`
	Alliance  string `json:"alliance" form:"alliance"`
}

// Search handles POST /api/search
func (h *Handlers) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Prevent same-city searches (e.g. LGA to EWR)
	same, err := h.DB.SameCity(req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if same {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end must be different cities (e.g. LGA and EWR are both NYC)"})
		return
	}
	args := flights.FlightSearch{
		Start:     req.Start,
		End:       req.End,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Cabin:     req.Cabin,
		Alliance:  req.Alliance,
	}
	result, err := flights.Explore(h.DB, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Cities returns list of known city codes
func (h *Handlers) Cities(c *gin.Context) {
	rows, err := h.DB.Query("SELECT DISTINCT city_code FROM iata_cities ORDER BY city_code")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var cities []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		cities = append(cities, code)
	}
	c.JSON(http.StatusOK, cities)
}
