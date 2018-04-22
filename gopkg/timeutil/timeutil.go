package timeutil

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.Scanner interface so it can be used as a scan destination, similar to sql.NullString
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// DefaultTime change datetimepicker.js date value to db standart time
func DefaultTime(d string, i ...int) string {
	layout := "01/02/2006"
	t, err := time.Parse(layout, d)
	if err != nil {
		return d
	}
	if len(i) == 1 {
		t = t.AddDate(0, 0, i[0])
	}

	return t.Format("2006-01-02")
}

// DayRange function
func DayRange(start time.Time, end time.Time) int {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC).Unix()
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).Unix()
	dayDiff := float64(endDay-startDay) / 86400

	if dayDiff > 1 {
		dayDiff = math.Ceil(dayDiff)
	} else {
		dayDiff = math.Floor(dayDiff)
	}

	return int(dayDiff)
}

// SetTimeLocationWIB function
func SetTimeLocationWIB(dateTime time.Time) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now(), err
	}

	return dateTime.In(location), nil
}

// TranslateMonths to translate months name from en to id and id to en
func TranslateMonths(date string) string {

	r := strings.NewReplacer(
		//en -> id
		"January", "Januari",
		"February", "Februari",
		"March", "Maret",
		"May", "Mei",
		"June", "Juni",
		"July", "Juli",
		"August", "Agustus",
		"October", "Oktober",
		"December", "Desember",
		"Aug", "Agu",
		"Oct", "Okt",
		"Dec", "Des",

		//id -> en
		"Januari", "January",
		"Februari", "February",
		"Maret", "March",
		"Mei", "May",
		"Juni", "June",
		"Juli", "July",
		"Agustus", "August",
		"Oktober", "October",
		"Desember", "December",
		"Agu", "Aug",
		"Okt", "Oct",
		"Des", "Dec",
	)

	return r.Replace(date)
}

// TranslateMonthsMustBahasa force translation to bahasa
func TranslateMonthsMustBahasa(date string) string {
	r := strings.NewReplacer(
		//en -> id
		"January", "Januari",
		"February", "Februari",
		"March", "Maret",
		"May", "Mei",
		"June", "Juni",
		"July", "Juli",
		"August", "Agustus",
		"October", "Oktober",
		"December", "Desember",
		"Aug", "Agu",
		"Oct", "Okt",
		"Dec", "Des",
	)

	return r.Replace(date)
}

// TranslateMonthsMustEnglish force translation to english
func TranslateMonthsMustEnglish(date string) string {
	r := strings.NewReplacer(
		//id -> en
		"Januari", "January",
		"Februari", "February",
		"Maret", "March",
		"Mei", "May",
		"Juni", "June",
		"Juli", "July",
		"Agustus", "August",
		"Oktober", "October",
		"Desember", "December",
		"Agu", "Aug",
		"Ags", "Aug",
		"Okt", "Oct",
		"Des", "Dec",
	)

	return r.Replace(date)
}

// IsLessThan30Days function
func IsLessThan30Days(t time.Time) bool {
	if time.Since(t).Hours() < 720 {
		return true
	}
	return false
}

// HumanizeTime function
func HumanizeTime(f interface{}, langs ...string) (int, string) {
	lang := "en" // default english
	if len(langs) > 0 && langs[0] != "" {
		lang = langs[0]
	}

	switch f.(type) {
	case time.Time:
		return humanizeTime(TimeToDuration(f.(time.Time)), lang)
	case int64:
		d, err := SecondsToDuration(f.(int64))
		if err != nil {
			return 0, ""
		}
		return humanizeTime(d, lang)
	default:
		return 0, ""
	}
}

// TimeToDuration function
func TimeToDuration(t time.Time, c ...time.Time) time.Duration {
	var timeAnchor time.Time
	if len(c) > 0 {
		timeAnchor = c[0]
	} else {
		timeAnchor = time.Now()
	}
	return t.Sub(timeAnchor)
}

// SecondsToDuration function
func SecondsToDuration(t int64) (time.Duration, error) {
	s := strconv.FormatInt(t, 10) + "s"
	d, err := time.ParseDuration(s)
	if err != nil {
		return d, err
	}

	return d, nil
}

func humanizeTime(d time.Duration, lang string) (int, string) {
	t, k := calculateTime(d)

	if lang == "id" {
		k = bahasaTimeType[k]
	}

	if lang == "en" {
		if t > 1 || t < -1 {
			return t, k + "s"
		}
	}

	return t, k
}

var bahasaTimeType = map[string]string{
	"day":    "hari",
	"hour":   "jam",
	"minute": "menit",
}

func calculateTime(d time.Duration, kind ...string) (int, string) {
	var timeType string
	if len(kind) > 0 {
		timeType = kind[0]
	} else {
		timeType = "day"
	}

	var rangeTime int
	switch timeType {
	case "day":
		rangeTime = int(d.Hours())
	case "hour":
		rangeTime = int(d.Hours())
	case "minute":
		rangeTime = int(d.Minutes())
	}

	if timeType == "day" && rangeTime >= 24 {
		return (rangeTime / 24), timeType
	} else if timeType == "day" && rangeTime < 24 {
		return calculateTime(d, "hour")
	} else if timeType == "hour" && rangeTime == 0 {
		return calculateTime(d, "minute")
	}

	return rangeTime, timeType
}

// FirstDateNextMonth function
func FirstDateNextMonth(currentTime time.Time) (time.Time, error) {
	var nextMonth string
	if month := currentTime.Month() + 1; month < 10 {
		nextMonth = fmt.Sprintf("0%d", month)
	} else {
		nextMonth = strconv.FormatInt(int64(month), 10)
	}

	return time.Parse("2006-01-02", fmt.Sprintf("%d-%s-01", currentTime.Year(), nextMonth))
}

// FirstDateThisMonth function
func FirstDateThisMonth(currentTime time.Time) (time.Time, error) {
	cYear, cMonth, _ := currentTime.Date()

	return time.Date(cYear, cMonth, 1, 0, 0, 0, 0, currentTime.Location()), nil
}

// LastHour function
func LastHour(currentTime time.Time) time.Time {
	return CalculateInterval(currentTime, 1, "hour", true)
}

// Yesterday function
func Yesterday(currentTime time.Time) time.Time {
	return CalculateInterval(currentTime, 1, "day", true)
}

// LastMonth function
func LastMonth(currentTime time.Time) time.Time {
	return CalculateInterval(currentTime, 1, "month", true)
}

// CalculateInterval function
func CalculateInterval(currentTime time.Time, interval int, unit string, backward bool) time.Time {
	sign := 1
	if backward {
		sign *= -1
	}

	switch unit {
	case "hour":
		return currentTime.Add(time.Duration(sign*interval) * time.Hour)
	case "day":
		return currentTime.AddDate(0, 0, sign*interval)
	case "month":
		return currentTime.AddDate(0, sign*interval, 0)
	}

	return currentTime
}

// JackOfAllDates supposed to convert ANY date layout
func JackOfAllDates(s string) (time.Time, error) {
	r := strings.NewReplacer(
		"-", " ",
		"+", " ",
		"/", " ",
	)
	result := r.Replace(s)

	var (
		newTime time.Time
		err     error
	)

	for _, dateFormat := range []string{
		"2 Jan 2006",
		"2006 01 02",
		"02 01 2006",
	} {
		if newTime, err = time.Parse(dateFormat, result); err == nil {
			return newTime, nil
		}
	}
	return newTime, fmt.Errorf("No format matches")
}
