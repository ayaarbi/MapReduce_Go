package mapreduce

import (
	"fmt"
	"project/common"
	"log"
)

func DoMap(jobName string, mapTaskNumber int, inFile string, nReduce int, mapF func(string, string) []common.KeyValue) {
	contents, err := common.ReadFile(inFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	kvs := mapF(inFile, contents)
	partitions := make([][]common.KeyValue, nReduce)
	for _, kv := range kvs {
		i := common.Ihash(kv.Key) % nReduce
		partitions[i] = append(partitions[i], kv)
	}
	for i := 0; i < nReduce; i++ {
		fileName := common.ReduceName(jobName, mapTaskNumber, i)
		common.WriteKeyValuesToFile(fileName, partitions[i])
	}
}

func DoReduce(jobName string, reduceTaskNumber int, nMap int, reduceF func(string, []string) string) {
	fmt.Printf("Starting reduce task %d with nMap = %d\n", reduceTaskNumber, nMap)

	var interFiles []string
	for i := 0; i < nMap; i++ {
		file := common.ReduceName(jobName, i, reduceTaskNumber)
		fmt.Println("Looking for intermediate file:", file)
		interFiles = append(interFiles, file)
	}

	results := common.MergeAndReduce(interFiles, reduceF)

	if len(results) == 0 {
		log.Printf("Reduce task %d: no results (check intermediate files exist!)", reduceTaskNumber)
	} else {
		log.Printf("Reduce task %d produced %d results", reduceTaskNumber, len(results))
	}

	outputFile := fmt.Sprintf("mrtmp.job-res-%d", reduceTaskNumber)
	fmt.Println("Writing final reduce output to:", outputFile)
	common.WriteKeyValuesToFile(outputFile, results)
}

