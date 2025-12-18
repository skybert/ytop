package pkg

type (
	SortKey int
)

func (sk SortKey) String() string {
	switch sk {
	case SortKeyCPU:
		return "cpu"
	case SortKeyName:
		return "name"
	case SortKeyMemory:
		return "memory"
	}
	return "unknown"
}

const (
	SortKeyCPU SortKey = iota
	SortKeyMemory
	SortKeyName
)

type Process struct {
	Pid     int
	Name    string
	Args    string
	RSS     uint64
	CPU     float64
	Env     []string
	Created int64
}

type YTopConf struct {
	HeaderForeground   string
	HeaderBackground   string
	SelectedForeground string
	SelectedBackground string
	SimpleView         bool
}
