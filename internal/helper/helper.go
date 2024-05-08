package helper

import (
	"fmt"
	"math"
	"strconv"
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

func PercentageCaluculation(moduleNumber int, values []int) string {
	total := moduleNumber * 100
	var current int
	for _, v := range values {
		current += v
	}
	p := (current * 100) / total
	percentage := math.Floor(float64(p))
	res := strconv.Itoa(int(percentage))

	if res == "100" {
		res = "completed"
	}

	return res
}
