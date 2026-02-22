package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddFlightRequest for POST /api/flights
type AddFlightRequest struct {
	Airline         string `json:"airline" binding:"required"`
	FlightNumber    string `json:"flight_number" binding:"required"`
	FromIATA        string `json:"from_iata" binding:"required"`
	ToIATA          string `json:"to_iata" binding:"required"`
	DepartureDate   string `json:"departure_date" binding:"required"`
	DepartureTime   string `json:"departure_time"`
	Confirmation    string `json:"confirmation"`
}

// ListFlights returns flights for the authenticated user
func (h *Handlers) ListFlights(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := h.DB.Query(`
		SELECT id, airline, flight_number, from_iata, to_iata, departure_date, departure_time, confirmation
		FROM booked_flights WHERE user_id = ? ORDER BY departure_date, departure_time
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var flights []gin.H
	for rows.Next() {
		var id int
		var airline, fn, from, to, date string
		var timeNull, confNull sql.NullString
		if err := rows.Scan(&id, &airline, &fn, &from, &to, &date, &timeNull, &confNull); err != nil {
			continue
		}
		time, conf := "", ""
		if timeNull.Valid {
			time = timeNull.String
		}
		if confNull.Valid {
			conf = confNull.String
		}
		flights = append(flights, gin.H{
			"id":               id,
			"airline":          airline,
			"flight_number":    fn,
			"from_iata":        from,
			"to_iata":          to,
			"departure_date":   date,
			"departure_time":   time,
			"confirmation":     conf,
		})
	}
	c.JSON(http.StatusOK, flights)
}

// AddFlight adds a flight for the authenticated user
func (h *Handlers) AddFlight(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req AddFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.DB.Exec(`
		INSERT INTO booked_flights (user_id, airline, flight_number, from_iata, to_iata, departure_date, departure_time, confirmation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, userID, req.Airline, req.FlightNumber, req.FromIATA, req.ToIATA, req.DepartureDate, req.DepartureTime, req.Confirmation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":               id,
		"airline":          req.Airline,
		"flight_number":    req.FlightNumber,
		"from_iata":        req.FromIATA,
		"to_iata":          req.ToIATA,
		"departure_date":   req.DepartureDate,
		"departure_time":   req.DepartureTime,
		"confirmation":     req.Confirmation,
	})
}

// DeleteFlight removes a flight (if owned by user)
func (h *Handlers) DeleteFlight(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := h.DB.Exec("DELETE FROM booked_flights WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
