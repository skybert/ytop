package internal

import (
	"fmt"

	mem "github.com/shirou/gopsutil/v4/mem"
	"skybert.net/ytop/pkg"
)

func MemorySummary(human bool) string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return ""
	}

	if human {
		return fmt.Sprintf(
			"total %s free %s",
			pkg.HumanBytes(v.Total),
			pkg.HumanBytes(v.Free),
		)
	} else {
		return fmt.Sprintf(
			"total %v free %v",
			v.Total,
			v.Free,
		)
	}

}
