package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"project/master"
	"project/worker"
	"project/mapreduce"
)

func main() {
	mode := flag.String("mode", "master", "Mode: master | worker")
	fileList := flag.String("files", "", "Comma-separated list of input files (for master)")
	flag.Parse()

	if *mode == "master" {
		if *fileList == "" {
			log.Fatal("Please provide input files with -files option")
		}
		files := strings.Split(*fileList, ",")
		m := &master.Master{
			Workers: make(map[string]*master.WorkerInfo),
		}
		for _, file := range files {

			m.Tasks = append(m.Tasks, master.Task{Type: "map", File: file, Status: "Idle"})
		}
		const nReduce = 3
		nMap := len(files)
		for i := 0; i < nReduce; i++ {
			m.Tasks = append(m.Tasks, master.Task{Type: "reduce", File: fmt.Sprintf("reduce-task-%d", i), Status: "Idle",NMap:   nMap})
		}
		go m.StartRPC()
		log.Println("Master RPC server started ")
		m.StartDashboard()
	} else if *mode == "worker" {
		w := &worker.Worker{
			Address: "localhost",
			MapF:    mapreduce.WordCountMapF,
			ReduceF: mapreduce.WordCountReduceF,
		}
		w.Run()
	} else {
		log.Fatal("Invalid mode. Use -mode=master or -mode=worker")
	}
}
