package main

import "sort"
import psutil "github.com/shirou/gopsutil/v4/process"

func getProcesses() []process {
	procs, err := psutil.Processes()
	if err != nil {
		return nil
	}

	var result []process

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

		result = append(result, process{
			Pid:  int(p.Pid),
			Name: name,
			Args: cmd,
			RSS:  mem.RSS / 1024,
			CPU:  cpu,
		})
	}

	return result
}

func (m *model) sortProcesses(processes []process) {
	sort.Slice(processes, func(i, j int) bool {
		pi := processes[i]
		pj := processes[j]

		switch m.sortKey {
		case SortKeyMemory:
			return pi.RSS > pj.RSS
		case SortKeyCPU:
			return pi.CPU > pj.CPU
		case SortKeyName:
			// Sort name ascending
			return pi.Name < pj.Name
		}
		return false
	})
}
