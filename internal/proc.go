package internal

import (
	psutil "github.com/shirou/gopsutil/v4/process"
	"skybert.net/ytop/pkg"
)

func Process(pid int) *pkg.Process {
	procs := Processes()
	for _, p := range procs {
		if p.Pid == pid {
			return &p
		}
	}

	return nil
}

// Processes returns the current running processes, represented using
// the ytop model.
func Processes() []pkg.Process {
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
		envVars, err := p.Environ()
		// Environ() isn't implemented yet on macOS, so just
		// ignore the err. For other errors, skip the process.
		if err != nil && err.Error() != "not implemented yet" {
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
			Env:  envVars,
		})
	}

	return result
}
