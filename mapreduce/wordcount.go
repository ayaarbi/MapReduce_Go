package mapreduce

import (
	"strconv"
	"strings"
	"unicode"
	"project/common"
)

func WordCountMapF(filename string, content string) []common.KeyValue {
	words := strings.FieldsFunc(content, func(r rune) bool { return !unicode.IsLetter(r) })
	kvs := []common.KeyValue{}
	for _, word := range words {
		kvs = append(kvs, common.KeyValue{Key: word, Value: "1"})
	}
	return kvs
}

func WordCountReduceF(key string, values []string) string {
	sum := 0
	for _, v := range values {
		n, _ := strconv.Atoi(v)
		sum += n
	}
	return strconv.Itoa(sum)
}
