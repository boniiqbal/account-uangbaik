package utils

import (
	"fmt"
	"time"
)

// SetToLateNight receive date parameter value "YYYY-MM-DD"
func SetToLateNight(date string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, parsed.Location()), nil
}

//CheckDateExpiredInput for validation
func CheckDateExpiredInput(enddate string) (bool, error) {

	tstart, err := SetToLateNight(time.Now().Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := SetToLateNight(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("Tanggal Expired tidak boleh kurang dari tanggal sekarang")
	}
	return true, err
}
