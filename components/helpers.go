package components

import (
	"strconv"

	"github.com/justestif/go-climbing/internal/database"
)

func gradeValue(n int) string {
	return strconv.Itoa(n)
}

func gradeName(n int) string {
	return "V" + strconv.Itoa(n)
}

func energyValue(log *database.SessionLog, def int) string {
	if log != nil && log.EnergyLevel.Valid {
		return strconv.Itoa(int(log.EnergyLevel.Int32))
	}
	return strconv.Itoa(def)
}

func sorenessValue(log *database.SessionLog, def int) string {
	if log != nil && log.Soreness.Valid {
		return strconv.Itoa(int(log.Soreness.Int32))
	}
	return strconv.Itoa(def)
}

func skinSelected(log *database.SessionLog, val string) bool {
	if log != nil && log.SkinCondition.Valid {
		return log.SkinCondition.String == val
	}
	return false
}
