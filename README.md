## To run the application

Here is an example of commands you could run
```
go run cmd/main.go -schedule="$(cat testdata/input.txt)" -offset=16:10
```

## To run tests

```
make test
```

cron.go holds the core business logic, its coverage only reaches 66.2% of statements. Cronjobs are difficult to test as they are time set functions, it would require long running tests. However, I have extensively tested nextTimeToRun function which holds the time computation to the next of our cronjobs.
```
cronscheduler/internal/cron     0.028s  coverage: 66.2% of statements
```
