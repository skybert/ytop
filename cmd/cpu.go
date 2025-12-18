package main

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
)

func CPUSummary() string {
	infos, err := cpu.Info()
	if err != nil || len(infos) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%v %v MHz",
		infos[0].ModelName,
		infos[0].Mhz)
}
