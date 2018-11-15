package crawlTime

import (
	"testing"
	"time"
)

func TestGetDate(t *testing.T) {
	t.Run("Should handle adding dates", func(t *testing.T) {
		now := time.Now()
		addedTime := GetDate(2)

		if addedTime.Day() != now.Day() + 2 {
			t.Error("Wrong day received")
		}
	})

	t.Run("Should handle subtracting dates", func(t *testing.T) {
		now := time.Now()
		addedTime := GetDate(-2)

		if addedTime.Day() != now.Day() - 2 {
			t.Error("Wrong day received")
		}
	})
}

func TestFindPreviousDate(t *testing.T) {
	dateFirst := "2018-08-12"
	datePrevious := "2018-8-11"

	actual := FindPreviousDate(dateFirst)
	if actual != datePrevious {
		t.Errorf("Expected: %v, got: %v", datePrevious, actual)
	}
}