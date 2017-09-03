package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	month31 = map[int]int{
		1:  1,
		3:  3,
		5:  5,
		7:  7,
		8:  8,
		10: 10,
		12: 12,
	}
)

/**
min hour day month weekday
*0-59 *0-23 *1-31 *1-12 *0-6
*/
func NearestFuture(clock string) (int64, error) {
	now := time.Now()
	timeMap := map[int]int{
		0: now.Minute(),
		1: now.Hour(),
		2: now.Day(),
		3: int(now.Month()),
		4: int(now.Weekday()),
	}
	digits := strings.Split(strings.Trim(clock, " "), " ")
	if len(digits) != 5 {
		return 0, fmt.Errorf("clock is wrong: %v", clock)
	}
	maxs := []int{59, 23, 31, 12, 6}
	mins := []int{0, 0, 1, 1, 0}
	seeds := map[string][]int{}
	nameMap := map[int]string{0: "minute", 1: "hour", 2: "day", 3: "month", 4: "weekday"}
	for i := 0; i < 5; i++ {
		ints := []int{}
		if strings.Index(digits[i], "*") >= 0 {
			ints = []int{mins[i], timeMap[i]}
			if timeMap[i] < maxs[i] {
				ints = append(ints, timeMap[i]+1)
			}
		} else {
			eles := strings.Split(digits[i], ",")
			for _, ele := range eles {
				tmp, err := strconv.ParseUint(ele, 10, 64)
				if err != nil {
					return 0, fmt.Errorf("clock is wrong: %v, %v", clock, err)
				}
				ints = append(ints, int(tmp))
			}
		}
		seeds[nameMap[i]] = ints
	}
	arrayT := []time.Time{}
	year := []int{now.Year(), now.Year() + 1}
	for _, y := range year {
		for _, m := range seeds["month"] {
			for _, d := range seeds["day"] {
				for _, h := range seeds["hour"] {
					for _, mi := range seeds["minute"] {
						if !isLeapYear(y) && m == 2 && d > 28 {
							continue
						}
						if m == 2 && d > 29 {
							continue
						}
						if !has31days(m) && d > 30 {
							continue
						}
						atime := time.Date(y, time.Month(m), d, h, mi, 0, 0, now.Location())
						arrayT = append(arrayT, atime)
					}
				}
			}
		}
	}
	for _, w := range seeds["weekday"] {
		days := w - timeMap[4]
		if days < 0 {
			days += 7
		}
		future := now.AddDate(0, 0, days)
		for _, h := range seeds["hour"] {
			for _, mi := range seeds["minute"] {
				atime := time.Date(future.Year(), future.Month(), future.Day(), h, mi, 0, 0, now.Location())
				arrayT = append(arrayT, atime)
			}
		}
	}
	minDuration := arrayT[0].Sub(now)
	for _, t := range arrayT {
		d := t.Sub(now)
		if d > minDuration {
			minDuration = d
		}
	}
	var x time.Time
	for _, t := range arrayT {
		// fmt.Println(t.String())
		d := t.Sub(now)
		if d > time.Nanosecond && d < minDuration {
			minDuration = d
			x = t
		}
	}
	fmt.Println("min is", x.String())
	//in case of less than 1 second, add 1
	return int64(minDuration.Seconds()) + 1, nil
}

func isLeapYear(year int) bool {
	if (year % 400) == 0 {
		return true
	}
	if (year%100) != 0 && (year%4) == 0 {
		return true
	}
	return false
}

func has31days(month int) bool {
	_, yes := month31[month]
	return yes
}
