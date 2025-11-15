package internal

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

func HumanUptime() string {
	d := Uptime()
	totalSeconds := int64(d.Seconds())

	days := totalSeconds / (24 * 3600)
	hours := (totalSeconds % (24 * 3600)) / 3600
	minutes := (totalSeconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d days, %02d:%02d", days, hours, minutes)
	}
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func Uptime() time.Duration {
	info, _ := host.Info()
	boot := time.Unix(int64(info.BootTime), 0)
	return time.Since(boot)
}
