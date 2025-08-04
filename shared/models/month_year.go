package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const monthYearLayout = "01-2006"

type MonthYear time.Time

func (my *MonthYear) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	if str == "" {
		return fmt.Errorf("date cannot be empty")
	}

	t, err := time.Parse(monthYearLayout, str)
	if err != nil {
		return fmt.Errorf("invalid date format, expected MM-YYYY")
	}

	*my = MonthYear(t)
	return nil
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(my).Format(monthYearLayout))
}

func (my MonthYear) Time() time.Time {
	return time.Time(my)
}

func (my MonthYear) Value() (driver.Value, error) {
	return time.Time(my), nil
}
