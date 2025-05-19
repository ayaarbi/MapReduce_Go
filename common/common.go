package common

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"sort"
)

type KeyValue struct {
	Key   string
	Value string
}

func ReduceName(jobName string, mapTask, reduceTask int) string {
	return fmt.Sprintf("mrtmp.%s-%d-%d", jobName, mapTask, reduceTask)
}

func Ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}

func EncodeJSON(kvs []KeyValue) ([]byte, error) {
	return json.Marshal(kvs)
}

func DecodeJSON(data []byte) ([]KeyValue, error) {
	var kvs []KeyValue
	err := json.Unmarshal(data, &kvs)
	return kvs, err
}

func ReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

func WriteKeyValuesToFile(filename string, kvs []KeyValue) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	for _, kv := range kvs {
		err := enc.Encode(&kv)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadKeyValuesFromFile(filename string) ([]KeyValue, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	var kvs []KeyValue
	for {
		var kv KeyValue
		if err := dec.Decode(&kv); err != nil {
			break
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
}

func MergeAndReduce(interFiles []string, reduceF func(string, []string) string) []KeyValue {
	intermediate := make(map[string][]string)
	for _, file := range interFiles {
		kvs, _ := ReadKeyValuesFromFile(file)
		for _, kv := range kvs {
			intermediate[kv.Key] = append(intermediate[kv.Key], kv.Value)
		}
	}

	var keys []string
	for k := range intermediate {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []KeyValue
	for _, k := range keys {
		results = append(results, KeyValue{k, reduceF(k, intermediate[k])})
	}
	return results
}
