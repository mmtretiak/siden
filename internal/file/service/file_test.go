package service

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mmtretiak/siden/pkg/logger"
)

const (
	pathToFile = "./file_example.txt"
	bufferSize = 1024 * 100
)

func TestService_ReadFile(t *testing.T) {
	log, err := logger.New(logger.ProductionLogType)
	if err != nil {
		t.Fatal(err)
	}

	service := New(bufferSize, log)

	res, err := service.ReadFile(pathToFile)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	words := strings.Split(string(data), "\n")
	uniqueWords := make(map[string]struct{})

	for _, word := range words {
		uniqueWords[word] = struct{}{}
	}

	if len(uniqueWords) != len(res) {
		t.Fatalf("invalid res, expected len: %d, actual: %d", len(uniqueWords), len(res))
	}
}

func BenchmarkService_ReadFile(b *testing.B) {
	log, err := logger.New(logger.ProductionLogType)
	if err != nil {
		b.Fatal(err)
	}

	service := New(bufferSize, log)

	for n := 0; n < b.N; n++ {
		_, err := service.ReadFile(pathToFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}
