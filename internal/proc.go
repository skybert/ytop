package internal

import (
	psutil "github.com/shirou/gopsutil/v4/process"
	"skybert.net/ytop/pkg"
)

func GetProcesses() []pkg.Process {
	procs, err := psutil.Processes()
	if err != nil {
		return nil
	}

	var result []pkg.Process

	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		args, _ := p.Cmdline()

		mem, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			continue
		}

		cmd := args
		if cmd == "" {
			cmd = name
		}

		result = append(result, pkg.Process{
			Pid:  int(p.Pid),
			Name: name,
			Args: cmd,
			RSS:  mem.RSS / 1024,
			CPU:  cpu,
		})
	}

	return result
}
