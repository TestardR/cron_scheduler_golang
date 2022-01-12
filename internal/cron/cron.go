package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	minute = iota
	hour
	cmd
)

var counter int

const (
	day        = 24
	timeFormat = "15:04"
)

// Job holds the state of a cronjob.
type Job struct {
	hour   string //  hour of the day
	minute string // minutes past the hour
	cmd    string //  command to run
}

// Offset holds the state of the given time command line argument (-offset).
type Offset struct {
	hour   int
	minute int
}

// Schedule function calls jobs according to the given job.
func (j Job) Schedule(o Offset) error {
	hour, minute, err := j.parseJobToHourAndMinute()
	if err != nil {
		return err
	}

	callback := func() {
		switch {
		case runEveryTime(j.hour) && runEveryTime(j.minute):
			everyMinute(1)
		case runEveryTime(j.hour):
			everyMinute(time.Duration(minute))
		case runEveryTime(j.minute):
			everyMinuteForOneHour(1)
		default:
			everyDay()
		}
	}

	time.AfterFunc(duration(hour, minute, o), callback)

	output := j.timeToNextRun(hour, minute, o)

	// log time not next run
	fmt.Println(output)

	return nil
}

func (j Job) timeToNextRun(hour, minute int, o Offset) string {
	t := time.Now()
	d := duration(hour, minute, o)
	offsetTime := time.Date(t.Year(), t.Month(), t.Day(), o.hour, o.minute, 0, 0, t.Location())
	c := offsetTime.Add(d)

	timer, _ := time.Parse(timeFormat, fmt.Sprintf("%d:%d", hour, minute))
	timeWithOffset := timer.Add(time.Hour*time.Duration(o.hour) + time.Minute*time.Duration(o.minute))

	switch {
	case runEveryTime(j.hour) && runEveryTime(j.minute) || runEveryTime(j.hour):
		output := fmt.Sprintf("%d:%d today - %s", timeWithOffset.Hour(), timeWithOffset.Minute(), j.cmd)

		return output
	default:
		day := "today"
		if c.Day() != offsetTime.Day() {
			day = "tomorrow"
		}

		output := fmt.Sprintf("%d:%d %s - %s", hour, minute, day, j.cmd)

		return output
	}
}

func (j Job) parseJobToHourAndMinute() (int, int, error) {
	var hour int

	var minute int

	if j.hour != "*" {
		h, err := strconv.Atoi(j.hour)
		if err != nil {
			return 0, 0, err
		}

		hour = h
	}

	if j.minute != "*" {
		m, err := strconv.Atoi(j.minute)
		if err != nil {
			return 0, 0, err
		}

		minute = m
	}

	return hour, minute, nil
}

// duration function computes the time difference between the offset time and the cronjob start.
func duration(hour, minute int, offset Offset) time.Duration {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, t.Location())
	o := time.Date(t.Year(), t.Month(), t.Day(), offset.hour, offset.minute, 0, 0, t.Location())

	if o.After(n) {
		n = n.Add(day * time.Hour)
	}

	d := n.Sub(o)

	return d
}

// ParseOffset function parses -offset command line argument.
func ParseOffset(arg string) (*Offset, error) {
	minute, err := strconv.Atoi(arg[3:])
	if err != nil {
		return nil, err
	}

	hour, err := strconv.Atoi(arg[:2])
	if err != nil {
		return nil, err
	}

	return &Offset{
		minute: minute,
		hour:   hour,
	}, nil
}

// ParseSchedule function parses -schedule command line argument.
func ParseSchedule(arg string) []Job {
	jobs := make([]Job, 0)

	for _, entry := range strings.Split(arg, "\n") {
		var c Job

		for i, v := range strings.Fields(entry) {
			switch i {
			case minute:
				c.minute = v
			case hour:
				c.hour = v
			case cmd:
				c.cmd = v
			}
		}

		jobs = append(jobs, c)
	}

	return jobs
}

func everyMinuteForOneHour(n time.Duration) {
	counter++

	if counter == 60 {
		time.Sleep(day * time.Hour)
	}

	time.AfterFunc(n*time.Minute, func() {
		everyMinuteForOneHour(1)
	})
}

func everyMinute(n time.Duration) {
	time.AfterFunc(n*time.Minute, func() {
		everyMinute(n)
	})
}

func everyDay() {
	time.AfterFunc(day*time.Hour, everyDay)
}

func runEveryTime(input string) bool {
	return input == "*"
}
