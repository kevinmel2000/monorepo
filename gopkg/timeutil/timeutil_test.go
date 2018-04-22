package timeutil_test

import (
	"testing"
	"time"

	"github.com/lab46/monorepo/gopkg/timeutil"
)

func TestTimeToDuration(t *testing.T) {
	minuteTest, _ := time.Parse(time.RFC3339, "2016-05-05T09:00:59Z")
	minuteTest2, _ := time.Parse(time.RFC3339, "2016-05-05T09:30:00Z")
	minuteTest3, _ := time.Parse(time.RFC3339, "2016-05-05T08:30:00Z")

	hourTest, _ := time.Parse(time.RFC3339, "2016-05-06T08:00:00Z")
	hourTest2, _ := time.Parse(time.RFC3339, "2016-05-05T10:00:00Z")
	hourTest3, _ := time.Parse(time.RFC3339, "2016-05-05T08:00:00Z")

	dayTest, _ := time.Parse(time.RFC3339, "2016-05-09T09:00:00Z")
	dayTest2, _ := time.Parse(time.RFC3339, "2016-05-06T09:30:00Z")

	anchorTest, _ := time.Parse(time.RFC3339, "2016-05-05T09:00:00Z")

	testcases := []struct {
		input  time.Time
		output time.Duration
	}{
		{minuteTest, 59 * time.Second},
		{minuteTest2, 30 * time.Minute},
		{minuteTest3, -30 * time.Minute},
		{hourTest, 23 * time.Hour},
		{hourTest2, 1 * time.Hour},
		{hourTest3, -1 * time.Hour},
		{dayTest, 4 * 24 * time.Hour},
		{dayTest2, (24 * time.Hour) + (30 * time.Minute)},
	}

	for _, v := range testcases {
		temp := timeutil.TimeToDuration(v.input, anchorTest)
		if temp != v.output {
			t.Errorf("Expected %v, got %v", v.output, temp)
		}
	}
}

func TestSecondsToDuration(t *testing.T) {
	testcases := []struct {
		input  int64
		output time.Duration
	}{
		{int64(0), 0 * time.Second},
		{int64(59), 59 * time.Second},
		{int64(-59), -59 * time.Second},
		{int64(6 * 60), 6 * time.Minute},
		{int64(60 * 60), 1 * time.Hour},
		{int64((60 * 60) + 1), (1 * time.Hour) + (1 * time.Second)},
	}

	for _, v := range testcases {
		temp, _ := timeutil.SecondsToDuration(v.input)
		if temp != v.output {
			t.Errorf("Expected %v, got %v", v.output, temp)
		}
	}
}

func TestHumanizeTime(t *testing.T) {
	testcases := []struct {
		input          interface{}
		outputTimeInt  int
		outputTimeType string
		lang           string
	}{
		// Test #1
		{int64(0), 0, "minute", ""},
		// Test #2
		{int64(1), 0, "minute", ""},
		// Test #3
		{int64(1 * 60), 1, "minute", ""},
		// Test #4
		{int64((1 * 60) + 59), 1, "minute", ""},
		// Test #5
		{int64((2 * 60) + 1), 2, "minutes", ""},
		// Test #6
		{int64(1 * 60 * 60), 1, "hour", ""},
		// Test #7
		{int64((1 * 60 * 60) + (3 * 60) + 39), 1, "hour", ""},
		// Test #8
		{int64((10 * 60 * 60) + (3 * 60) + 39), 10, "hours", ""},
		// Test #9
		{int64((23 * 60 * 60) + (59 * 60) + 59), 23, "hours", ""},
		// Test #10
		{int64(1 * 24 * 60 * 60), 1, "day", ""},
		// Test #11
		{int64(20 * 24 * 60 * 60), 20, "days", ""},
		// Test #12
		{int64((20 * 24 * 60 * 60) + (60 * 60)), 20, "days", ""},
		// Test #13
		{int64(1*60*60) * int64(time.Now().Sub(time.Now().AddDate(0, 0, 2)).Hours()) * -1, 2, "hari", "id"},
		// Test #14
		{int64(time.Now().Sub(time.Now().Add(time.Second*72000)).Minutes()) * -1, 20, "menit", "id"},
	}

	for k, v := range testcases {
		expiryTime, timeType := timeutil.HumanizeTime(v.input, v.lang)
		if expiryTime != v.outputTimeInt || timeType != v.outputTimeType {
			t.Errorf("Test #%d, expected %v and %v, got %v and %v", k, v.outputTimeInt, v.outputTimeType, expiryTime, timeType)
		}
	}
}
