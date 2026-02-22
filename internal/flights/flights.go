package flights

import (
	"strings"
	"triangle_travel/internal/db"
)

// FlightSearch represents search parameters
type FlightSearch struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Cabin     string `json:"cabin"`
	Alliance  string `json:"alliance"`
}

// Normalize ensures uppercase and defaults
func (f *FlightSearch) Normalize() {
	f.Start = strings.ToUpper(strings.TrimSpace(f.Start))
	f.End = strings.ToUpper(strings.TrimSpace(f.End))
	if f.Cabin == "" {
		f.Cabin = "economy"
	}
	if f.Alliance == "" {
		f.Alliance = "None"
	}
}

// TriangleResult holds places to explore
type TriangleResult struct {
	DriveThenFly map[string]float64 `json:"driveThenFly"` // IATA -> distance or price delta
	FlyThenFly   map[string]float64 `json:"flyThenFly"`   // IATA -> price delta
	AvgPrice     float64            `json:"avgPrice"`     // placeholder, -1 if unknown
}

// Explore returns triangle travel options from the database
func Explore(database *db.DB, args FlightSearch) (*TriangleResult, error) {
	args.Normalize()

	result := &TriangleResult{
		DriveThenFly: make(map[string]float64),
		FlyThenFly:   make(map[string]float64),
		AvgPrice:     -1,
	}

	// Distances: places you can drive/train to then fly
	distances, err := database.GetDistancesFrom(args.End)
	if err != nil {
		return result, err
	}
	for iata, dist := range distances {
		if dist >= 55 && dist <= 300 {
			result.DriveThenFly[iata] = dist
		}
	}

	// City routes: places you can fly to then fly out of
	routes, err := database.GetRoutesFromWithFallback(args.End, args.Alliance)
	if err != nil {
		return result, err
	}
	for _, iata := range routes {
		result.FlyThenFly[iata] = 0 // DB doesn't have price; frontend can link to Kayak
	}

	return result, nil
}
