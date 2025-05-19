package tests

import (
	"net/rpc"
	"os"
	"project/common"
	"project/mapreduce"
	"project/master"
	"testing"
	"time"
)

func TestWordCountMapF(t *testing.T) {
	input := "Go is expressive, concise, clean, and efficient."
	expected := map[string]int{
		"Go": 1, "is": 1, "expressive": 1, "concise": 1, "clean": 1, "and": 1, "efficient": 1,
	}
	kvs := mapreduce.WordCountMapF("", input)
	result := make(map[string]int)
	for _, kv := range kvs {
		result[kv.Key]++
	}
	for word, count := range expected {
		if result[word] != count {
			t.Errorf("Expected %d for word '%s', got %d", count, word, result[word])
		}
	}
}

func TestWordCountReduceF(t *testing.T) {
	key := "test"
	values := []string{"1", "1", "1"}
	expected := "3"
	result := mapreduce.WordCountReduceF(key, values)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestWriteAndReadKeyValues(t *testing.T) {
	file := "test_kv.json"
	kvs := []common.KeyValue{
		{Key: "hello", Value: "1"},
		{Key: "world", Value: "2"},
	}
	err := common.WriteKeyValuesToFile(file, kvs)
	if err != nil {
		t.Fatalf("Error writing file: %v", err)
	}
	readKVs, err := common.ReadKeyValuesFromFile(file)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	if len(kvs) != len(readKVs) {
		t.Fatalf("Mismatch in KV length: got %d, want %d", len(readKVs), len(kvs))
	}
	os.Remove(file)
}

func TestMasterWorkerRPC(t *testing.T) {
	m := &master.Master{
		Workers: make(map[string]*master.WorkerInfo),
	}
	m.Tasks = append(m.Tasks, master.Task{Type: "map", File: "test.txt", Status: "Idle"})
	m.Tasks = append(m.Tasks, master.Task{Type: "reduce", File: "reduce-0", Status: "Idle"})

	// Start RPC server
	go m.StartRPC()
	time.Sleep(1 * time.Second) // allow server to start

	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		t.Fatalf("Failed to dial RPC: %v", err)
	}
	defer client.Close()

	reply := &struct {
		Task  master.Task
		Index int
	}{}
	err = client.Call("Master.GetTask", &struct{}{}, reply)
	if err != nil {
		t.Errorf("RPC GetTask failed: %v", err)
	}

	if reply.Task.Status != "InProgress" {
		t.Errorf("Expected task to be InProgress, got %s", reply.Task.Status)
	}

	doneArgs := &master.DoneArgs{Index: reply.Index, Worker: "testWorker"}
	err = client.Call("Master.ReportTaskDone", doneArgs, &struct{}{})
	if err != nil {
		t.Errorf("RPC ReportTaskDone failed: %v", err)
	}

	if m.Tasks[reply.Index].Status != "Done" {
		t.Errorf("Expected task to be Done, got %s", m.Tasks[reply.Index].Status)
	}
}
