package service

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mmtretiak/siden/pkg/logger"
)

const (
	fileName   = "file_example.txt"
	bufferSize = 1024
)

func TestService_ReadFile(t *testing.T) {
	log, err := logger.New(logger.ProductionLogType)
	if err != nil {
		t.Fatal(err)
	}

	service := New(bufferSize, "./", log)

	res, err := service.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	words := strings.Split(string(data), "\n")
	uniqueWords := make(map[string]struct{})

	for _, word := range words {
		if word == "" {
			continue
		}

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

	service := New(bufferSize, "./", log)

	for n := 0; n < b.N; n++ {
		_, err := service.ReadFile(fileName)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportAllocs()
	}
}
