package main

import (
	"sort"

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

func SortProcesses(processes []pkg.Process, sortKey pkg.SortKey) {
	sort.Slice(processes, func(i, j int) bool {
		pi := processes[i]
		pj := processes[j]

		switch sortKey {
		case pkg.SortKeyMemory:
			return pi.RSS > pj.RSS
		case pkg.SortKeyCPU:
			return pi.CPU > pj.CPU
		case pkg.SortKeyName:
			// Sort name ascending
			return pi.Name < pj.Name
		}
		return false
	})
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

		created, err := p.CreateTime()
		if err != nil {
			continue
		}

		result = append(result, pkg.Process{
			Pid:     int(p.Pid),
			Name:    name,
			Args:    cmd,
			RSS:     mem.RSS,
			CPU:     cpu,
			Env:     envVars,
			Created: created,
		})
	}

	return result
}
