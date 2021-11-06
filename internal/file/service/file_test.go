package service

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mmtretiak/siden/pkg/logger"
	"github.com/stretchr/testify/assert"
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

	uniqueWords = make(map[string]struct{})

	for _, word := range res {
		if _, ok := uniqueWords[word]; ok {
			t.Fatalf("found duplication of word: %s", word)
		}

		uniqueWords[word] = struct{}{}
	}
}

func TestService_WriteToFile(t *testing.T) {
	log, err := logger.New(logger.ProductionLogType)
	if err != nil {
		t.Fatal(err)
	}

	service := New(bufferSize, "./", log)

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(fileContent)

	testFileName := "test_" + fileName

	if err := service.WriteToFile(testFileName, buf); err != nil {
		t.Fatal(err)
	}

	writtenFile, err := os.Open(testFileName)
	if err != nil {
		t.Fatal(err)
	}

	writtenFileContent, err := ioutil.ReadAll(writtenFile)
	if err != nil {
		writtenFile.Close()
		t.Fatal(err)
	}

	writtenFile.Close()

	if err := os.Remove(testFileName); err != nil {
		panic(err)
	}

	assert.Equal(t, fileContent, writtenFileContent)
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
