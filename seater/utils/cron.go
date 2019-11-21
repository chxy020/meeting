package utils

import (
	"fmt"
	"time"
)

// ParseCronDuration parses cron spec and return duration time
func ParseCronDuration(cronExpression string) (d time.Duration) {
	var i int64
	_, err := fmt.Sscanf(cronExpression, "@every %ds", &i)
	if err != nil {
		return 0
	}
	d = time.Duration(i) * time.Second
	return
}

// get expiration time, 2 * duration or duration + 10s, returns the larger one
func getExpirationTime(d time.Duration) (expiration time.Duration) {
	a := time.Duration(2) * d
	b := d + 20*time.Second

	if a < b {
		return b
	}
	return a
}

// StatCron interval
var (
	MinerStatDuration   = ParseCronDuration("@every 10s")
	ClusterStatDuration = ParseCronDuration("@every 10s")
)

// Stat expiration time, set the expiration to 2 * cron interval, cause cron job may run with lantency
var (
	MinerStatExpiration   = getExpirationTime(MinerStatDuration)
	ClusterStatExpiration = getExpirationTime(ClusterStatDuration)
)
