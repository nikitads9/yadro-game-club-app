package process

import "math"

// calculateRevenue Функция вычисления выручки с округлением по часам
func calculateRevenue(pc *computer) int64 {
	dur := pc.SessionEnd.Sub(pc.SessionStart)
	return int64(math.Ceil(dur.Hours())) * price
}
