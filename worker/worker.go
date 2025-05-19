package worker

import (
	"log"
	"net/rpc"
	"project/common"
	"project/mapreduce"
	"project/master"
	"time"
)

type Worker struct {
	Address string
	MapF    func(string, string) []common.KeyValue
	ReduceF func(string, []string) string
}

const nReduce = 3

func (w *Worker) Run() {
	for {
		client, err := rpc.Dial("tcp", "localhost:1234")
		if err != nil {
			log.Println("Failed to connect to master:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		reply := &struct {
			Task  master.Task
			Index int
		}{}
		err = client.Call("Master.GetTask", &struct{}{}, reply)
		if err != nil {
			log.Println("GetTask RPC failed:", err)
			client.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		if reply.Task.File == "" {
			client.Close()
			time.Sleep(time.Second)
			continue
		}

		// Simulate crash or slow
		if !Simulate() {
			client.Close()
			return
		}

		if reply.Task.Type == "map" {
			mapreduce.DoMap("job", reply.Index, reply.Task.File, nReduce, w.MapF)
		} else {
    reduceIndex := reply.Index - reply.Task.NMap
    mapreduce.DoReduce("job", reduceIndex, reply.Task.NMap, w.ReduceF)
}


		doneArgs := &master.DoneArgs{Index: reply.Index, Worker: w.Address}
		err = client.Call("Master.ReportTaskDone", doneArgs, &struct{}{})
		if err != nil {
			log.Println("ReportTaskDone RPC failed:", err)
		}
		client.Close()
	}
}
