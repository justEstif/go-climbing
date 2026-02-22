package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/justestif/go-climbing/internal/database"
)

type WorkoutSet struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GradeRange  string `json:"grade_range,omitempty"`
	Duration    string `json:"duration,omitempty"`
	Sets        string `json:"sets,omitempty"`
}

type SessionPlan struct {
	Warmup  []WorkoutSet
	Main    []WorkoutSet
	Project []WorkoutSet
}

func GenerateFirstSession(userID, currentGrade, goalGrade int32, weaknesses []string, date time.Time) (database.CreateSessionParams, error) {
	plan := buildPlan(currentGrade, goalGrade, weaknesses)
	focusType := chooseFocusType(weaknesses)

	warmupJSON, err := json.Marshal(plan.Warmup)
	if err != nil {
		return database.CreateSessionParams{}, fmt.Errorf("marshal warmup: %w", err)
	}
	mainJSON, err := json.Marshal(plan.Main)
	if err != nil {
		return database.CreateSessionParams{}, fmt.Errorf("marshal main: %w", err)
	}
	projectJSON, err := json.Marshal(plan.Project)
	if err != nil {
		return database.CreateSessionParams{}, fmt.Errorf("marshal project: %w", err)
	}

	pgDate := pgtype.Date{}
	if err := pgDate.Scan(date); err != nil {
		return database.CreateSessionParams{}, fmt.Errorf("parse date: %w", err)
	}

	return database.CreateSessionParams{
		UserID:         pgtype.Int4{Int32: userID, Valid: true},
		SessionNumber:  1,
		Date:           pgDate,
		FocusType:      focusType,
		PlannedWarmup:  warmupJSON,
		PlannedMain:    mainJSON,
		PlannedProject: projectJSON,
	}, nil
}

func buildPlan(currentGrade, goalGrade int32, weaknesses []string) SessionPlan {
	easyGrade := currentGrade - 4
	if easyGrade < 0 {
		easyGrade = 0
	}
	projectGrade := currentGrade + 1
	if projectGrade > goalGrade {
		projectGrade = goalGrade
	}

	warmup := []WorkoutSet{
		{
			Name:        "Easy Traversals",
			Description: "Traverse the wall sideways on easy holds to warm up fingers and shoulders.",
			GradeRange:  fmt.Sprintf("V%d", easyGrade),
			Duration:    "10 min",
		},
		{
			Name:        "Mobility Stretching",
			Description: "Hip flexor stretches, shoulder circles, and wrist rotations.",
			Duration:    "5 min",
		},
	}

	var main []WorkoutSet
	focus := firstWeakness(weaknesses)
	switch focus {
	case "footwork":
		main = []WorkoutSet{
			{
				Name:        "Slab Drills",
				Description: "Climb slab routes focusing on precise foot placement. No hands allowed on warm-up laps.",
				GradeRange:  fmt.Sprintf("V%d–V%d", currentGrade-2, currentGrade),
				Sets:        "4 laps",
			},
			{
				Name:        "Silent Feet",
				Description: "Climb routes placing your feet silently — no scraping. Focus on reading foot holds before stepping.",
				GradeRange:  fmt.Sprintf("V%d", currentGrade),
				Sets:        "3 routes",
			},
		}
	case "power":
		main = []WorkoutSet{
			{
				Name:        "Limit Bouldering",
				Description: "Attempt problems at or slightly above your current limit. Rest 3–5 min between attempts.",
				GradeRange:  fmt.Sprintf("V%d–V%d", currentGrade, currentGrade+1),
				Sets:        "5 attempts per problem, 3 problems",
			},
			{
				Name:        "Campus Ladder",
				Description: "Campus board ladder rungs 1-3-5. Focus on explosive pull.",
				Sets:        "3 sets, rest 3 min",
			},
		}
	case "endurance":
		main = []WorkoutSet{
			{
				Name:        "4x4s",
				Description: "Choose 4 moderate routes. Climb each 4 times in a row without rest, then rest 3 min. Repeat.",
				GradeRange:  fmt.Sprintf("V%d–V%d", currentGrade-2, currentGrade-1),
				Sets:        "4 rounds",
			},
			{
				Name:        "Continuous Climbing",
				Description: "Climb without stopping for 20 minutes on easy terrain.",
				GradeRange:  fmt.Sprintf("V%d", easyGrade),
				Duration:    "20 min",
			},
		}
	default:
		main = []WorkoutSet{
			{
				Name:        "Volume Pyramid",
				Description: "Climb progressively harder routes: 4 at base grade, 3 one grade up, 2 two grades up, 1 three grades up.",
				GradeRange:  fmt.Sprintf("V%d–V%d", currentGrade-2, currentGrade+1),
				Sets:        "10 routes total",
			},
		}
	}

	project := []WorkoutSet{
		{
			Name:        fmt.Sprintf("Project Attempts (V%d)", projectGrade),
			Description: fmt.Sprintf("Work your project at V%d. Break it into sections, identify the crux, and try different beta.", projectGrade),
			GradeRange:  fmt.Sprintf("V%d", projectGrade),
			Sets:        "4–5 attempts",
		},
	}

	return SessionPlan{
		Warmup:  warmup,
		Main:    main,
		Project: project,
	}
}

func firstWeakness(weaknesses []string) string {
	if len(weaknesses) == 0 {
		return ""
	}
	return weaknesses[0]
}

func chooseFocusType(weaknesses []string) string {
	w := firstWeakness(weaknesses)
	switch w {
	case "footwork", "technique":
		return "technique"
	case "power":
		return "strength"
	case "endurance":
		return "endurance"
	default:
		return "general"
	}
}

func DecodeSessionPlan(s database.Session) (SessionPlan, error) {
	var warmup, main, project []WorkoutSet

	if err := json.Unmarshal(s.PlannedWarmup, &warmup); err != nil {
		return SessionPlan{}, fmt.Errorf("decode warmup: %w", err)
	}
	if err := json.Unmarshal(s.PlannedMain, &main); err != nil {
		return SessionPlan{}, fmt.Errorf("decode main: %w", err)
	}
	if err := json.Unmarshal(s.PlannedProject, &project); err != nil {
		return SessionPlan{}, fmt.Errorf("decode project: %w", err)
	}

	return SessionPlan{
		Warmup:  warmup,
		Main:    main,
		Project: project,
	}, nil
}

func ParseWeaknesses(raw []byte) ([]string, error) {
	var weaknesses []string
	if err := json.Unmarshal(raw, &weaknesses); err != nil {
		return nil, err
	}
	return weaknesses, nil
}

func MarshalWeaknesses(weaknesses []string) ([]byte, error) {
	if len(weaknesses) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(weaknesses)
}

type LoggedRoute struct {
	Grade int    `json:"grade"`
	Style string `json:"style"`
	Count int    `json:"count"`
}

func EncodeRoutesLogged(routes []LoggedRoute) ([]byte, error) {
	if len(routes) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(routes)
}

func DecodeRoutesLogged(data []byte) ([]LoggedRoute, error) {
	var routes []LoggedRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return nil, err
	}
	return routes, nil
}
