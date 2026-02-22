package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/components"
	"github.com/justestif/go-climbing/internal/database"
	"github.com/justestif/go-climbing/internal/middleware"
)

type GradePoint struct {
	Date  string `json:"x"`
	Grade int    `json:"y"`
}

type WellnessPoint struct {
	Date     string `json:"date"`
	Energy   *int   `json:"energy"`
	Soreness *int   `json:"soreness"`
}

type SessionFrequencyPoint struct {
	Date      string `json:"x"`
	DayOfWeek int    `json:"y"`
	Count     int    `json:"v"`
}

func ProgressPage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.SessionManager.GetInt(r.Context(), "userID")

	user, err := database.DB.GetUser(r.Context(), int32(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	logs, err := database.DB.ListSessionLogsByUser(r.Context(), pgtype.Int4{Int32: int32(userID), Valid: true})
	if err != nil {
		logs = nil
	}

	stats := components.ProgressStats{
		TotalLogged:  len(logs),
		CurrentGrade: user.CurrentMaxGrade,
		GoalGrade:    user.GoalGrade,
	}

	// Reverse logs so oldest-first for chart display (DB returns DESC)
	gradePoints := make([]GradePoint, 0)
	wellnessPoints := make([]WellnessPoint, 0)
	sessionCounts := make(map[string]int)
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]
		date := log.LoggedAt.Time.Format("2006-01-02")

		sessionCounts[date]++

		if log.NewMaxGrade.Valid {
			gradePoints = append(gradePoints, GradePoint{
				Date:  date,
				Grade: int(log.NewMaxGrade.Int32),
			})
		}

		if log.EnergyLevel.Valid || log.Soreness.Valid {
			wp := WellnessPoint{Date: date}
			if log.EnergyLevel.Valid {
				v := int(log.EnergyLevel.Int32)
				wp.Energy = &v
			}
			if log.Soreness.Valid {
				v := int(log.Soreness.Int32)
				wp.Soreness = &v
			}
			wellnessPoints = append(wellnessPoints, wp)
		}
	}

	frequencyPoints := make([]SessionFrequencyPoint, 0, len(sessionCounts))
	for date, count := range sessionCounts {
		t, _ := time.Parse("2006-01-02", date)
		// Convert Go weekday (0=Sun..6=Sat) to ISO weekday (1=Mon..7=Sun)
		w := int(t.Weekday())
		if w == 0 {
			w = 7
		}
		frequencyPoints = append(frequencyPoints, SessionFrequencyPoint{
			Date:      date,
			DayOfWeek: w,
			Count:     count,
		})
	}

	gradeJSON, err := json.Marshal(gradePoints)
	if err != nil {
		gradeJSON = []byte("[]")
	}
	wellnessJSON, err := json.Marshal(wellnessPoints)
	if err != nil {
		wellnessJSON = []byte("[]")
	}
	frequencyJSON, err := json.Marshal(frequencyPoints)
	if err != nil {
		frequencyJSON = []byte("[]")
	}

	components.ProgressPage(stats, string(gradeJSON), string(wellnessJSON), string(frequencyJSON)).Render(r.Context(), w)
}
