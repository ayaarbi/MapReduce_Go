package master

type State struct {
	Tasks   []Task
	Workers map[string]*WorkerInfo
	Done    int
	Total   int
}

func (m *Master) Snapshot() State {
	m.mu.Lock()
	defer m.mu.Unlock()
	return State{
		Tasks:   m.Tasks,
		Workers: m.Workers,
		Done:    m.Done,
		Total:   len(m.Tasks),
	}
}