package format

import (
	"fmt"
	"strings"
	"time"
)

// Event универсальная функция для вывода исходящих событий
func Event(moment time.Time, eventID int, message ...any) {
	fmt.Printf("%s %d %v\n", moment.Format("15:04"), eventID, strings.Trim(fmt.Sprintf("%v", message), "[]"))
}

// FmtDuration Функция форматирования длительности сеанса
func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
