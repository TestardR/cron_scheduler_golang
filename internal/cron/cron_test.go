package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	t.Run("should_be_return_the_right_duration", func(t *testing.T) {
		o := Offset{
			hour:   15,
			minute: 00,
		}
		d := duration(19, 00, o)

		assert.Equal(t, float64(4), d.Hours())
	})

	t.Run("should_be_add_24hours_to_given_time_if_offset_is_before", func(t *testing.T) {
		o := Offset{
			hour:   19,
			minute: 00,
		}

		d := duration(15, 00, o)

		assert.Equal(t, float64(20), d.Hours())
	})
}

func TestParseOffset(t *testing.T) {
	t.Run("should_parse_offset", func(t *testing.T) {
		o, err := ParseOffset("16:10")
		require.NoError(t, err)

		assert.Equal(t, 16, o.hour)
		assert.Equal(t, 10, o.minute)
	})

	t.Run("should_return_an_error_input_does_not_have_the_right_format", func(t *testing.T) {
		o, err := ParseOffset("test:test")
		require.Error(t, err)

		assert.Nil(t, o)
	})
}

func TestParseSchedule(t *testing.T) {
	t.Run("should_parse_schedule", func(t *testing.T) {
		const testCase = `30 1 /bin/run_me_daily
		45 * /bin/run_me_hourly`

		jobs := ParseSchedule(testCase)

		assert.Equal(t, "/bin/run_me_daily", jobs[0].cmd)
		assert.Equal(t, "30", jobs[0].minute)
		assert.Equal(t, "1", jobs[0].hour)
		assert.Equal(t, "/bin/run_me_hourly", jobs[1].cmd)
		assert.Equal(t, "45", jobs[1].minute)
		assert.Equal(t, "*", jobs[1].hour)
	})
}

func TestParseJobToHourAndMinute(t *testing.T) {
	t.Run("should_parse_job_to_hour_and_minute_if_*_is_given", func(t *testing.T) {
		job := Job{
			hour:   "*",
			minute: "*",
		}

		hour, minute, err := job.parseJobToHourAndMinute()
		require.NoError(t, err)

		assert.Equal(t, 0, hour)
		assert.Equal(t, 0, minute)
	})

	t.Run("should_parse_job_to_hour_and_minute", func(t *testing.T) {
		job := Job{
			hour:   "10",
			minute: "22",
		}

		hour, minute, err := job.parseJobToHourAndMinute()
		require.NoError(t, err)

		assert.Equal(t, 10, hour)
		assert.Equal(t, 22, minute)
	})
}

func TestTimeToNextRun(t *testing.T) {
	t.Run("should_output_the_right_message_case_1", func(t *testing.T) {
		job := Job{
			hour:   "1",
			minute: "30",
			cmd:    "/bin/run_me_daily",
		}
		o := Offset{hour: 16, minute: 10}

		output := job.timeToNextRun(1, 30, o)

		assert.Equal(t, "1:30 tomorrow - /bin/run_me_daily", output)
	})

	t.Run("should_output_the_right_message_case_2", func(t *testing.T) {
		job := Job{
			hour:   "*",
			minute: "45",
			cmd:    "/bin/run_me_hourly",
		}
		o := Offset{hour: 16, minute: 10}

		output := job.timeToNextRun(0, 45, o)

		assert.Equal(t, "16:55 today - /bin/run_me_hourly", output)
	})

	t.Run("should_output_the_right_message_case_3", func(t *testing.T) {
		job := Job{
			hour:   "*",
			minute: "*",
			cmd:    "/bin/run_me_every_minute",
		}
		o := Offset{hour: 16, minute: 10}

		output := job.timeToNextRun(0, 0, o)

		assert.Equal(t, "16:10 today - /bin/run_me_every_minute", output)
	})

	t.Run("should_output_the_right_message_case_4", func(t *testing.T) {
		job := Job{
			hour:   "19",
			minute: "00",
			cmd:    "/bin/run_me_sixty_times",
		}
		o := Offset{hour: 16, minute: 10}

		output := job.timeToNextRun(19, 0, o)

		assert.Equal(t, "19:0 today - /bin/run_me_sixty_times", output)
	})
}
