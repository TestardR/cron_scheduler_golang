package main

import (
	"cronscheduler/internal/cron"
	"cronscheduler/internal/logger"
	"errors"
	"flag"
	"os"
)

const (
	appName       = "cronscheduler"
	scheduleUsage = `Please enter a schedule in the following format:
30 1 /bin/run_me_daily
45 * /bin/run_me_hourly
* * /bin/run_me_every_minute
* 19 /bin/run_me_sixty_times
`
	offsetUsage = "Please enter a offset value in the following format: hh:mm"
)

var (
	errMissingArgument = errors.New("missing command line argument")
	errOffsetFormat    = errors.New("offset argument should respect the format HH:MM")
)

func main() {
	log := logger.New(appName)

	s := flag.String("schedule", "", scheduleUsage)
	t := flag.String("offset", "", offsetUsage)
	flag.Parse()

	if t == nil || s == nil {
		log.Fatal(errMissingArgument)
		os.Exit(1)
	}

	if len(*t) != 5 {
		log.Fatal(errOffsetFormat)
		os.Exit(1)
	}

	jobs := cron.ParseSchedule(*s)

	offset, err := cron.ParseOffset(*t)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Info("scheduling cronjobs started")

	for _, j := range jobs {
		err := j.Schedule(*offset)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("scheduling cronjobs finished")

	select {}
}
