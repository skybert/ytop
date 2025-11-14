package main

import "sort"

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
