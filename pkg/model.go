package pkg

type Process struct {
	Pid  int
	Name string
	Args string
	RSS  uint64
	CPU  float64
	Env  []string
}
