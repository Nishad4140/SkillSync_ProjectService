package helper

import (
	"fmt"
	"time"
)

func ConvertStringToDate(data string) (time.Time, error) {
	layOut := "02-01-2006"
	date, err := time.Parse(layOut, data)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while converting to time")
	}
	return date, nil
}
