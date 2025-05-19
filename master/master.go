package master

import (
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Task struct {
	Type   string // "map" or "reduce"
	File   string
	Status string // Idle, InProgress, Done
	Worker string
	NMap  int
}

type WorkerInfo struct {
	Address string
	Status  string `json:"Status"`
}

type Master struct {
	mu        sync.Mutex
	Tasks     []Task
	Workers   map[string]*WorkerInfo
	TaskIndex int
	Done      int
}

type Args struct{}

type TaskReply struct {
	Task Task
	Index int
}

type DoneArgs struct {
	Index  int
	Worker string
}

func (m *Master) GetTask(args *Args, reply *TaskReply) error {
	
	workerIP := "localhost"
	
	m.mu.Lock()
	defer m.mu.Unlock()
if _, exists := m.Workers[workerIP]; !exists {
		m.Workers[workerIP] = &WorkerInfo{Address: workerIP, Status: "Idle"}
	}

	for i := range m.Tasks {
		t := &m.Tasks[i]
		if t.Status == "Idle" {
			t.Status = "InProgress"
			t.Worker = workerIP
			m.Workers[workerIP].Status = "Busy"
			reply.Task = *t
			reply.Index = i
			go m.monitorTimeout(i)
			return nil
		}
	}
	return nil
}

func (m *Master) ReportTaskDone(args *DoneArgs, reply *Args) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if args.Index < len(m.Tasks) {
		m.Tasks[args.Index].Status = "Done"
		m.Tasks[args.Index].Worker = args.Worker
		m.Done++
	}
	return nil
}

func (m *Master) monitorTimeout(index int) {
	time.Sleep(10 * time.Second)
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Tasks[index].Status == "InProgress" {
		m.Tasks[index].Status = "Idle"
		m.Tasks[index].Worker = ""
	}
}

func (m *Master) StartRPC() {
	rpc.Register(m)
	l, _ := net.Listen("tcp", ":1234")
	go rpc.Accept(l)
}

