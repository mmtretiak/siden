package service

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/mmtretiak/siden/pkg/logger"
)

func New(bufferSize int, log logger.Logger) *service {
	return &service{
		bufferSize: bufferSize,
		log:        log,
	}
}

type service struct {
	bufferSize int
	log        logger.Logger
}

func (s *service) WriteToFile(pathToFile string, data []byte) error {
	file, err := os.Open(pathToFile)
	if err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s pathToFile=%s", "WriteToFile", "open file", err.Error(), pathToFile)
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s pathToFile=%s writtenSize=%v dataSize=%v", "WriteToFile",
			"write data into file", err.Error(), pathToFile, n, len(data))
		return err
	}

	return nil
}

func (s *service) ReadFile(pathToFile string) ([]string, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s pathToFile=%s", "ReadFile", "open file", err.Error(), pathToFile)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	uniqueWords := make(map[string]struct{})

	numberOfIterations := 0
	for {
		numberOfIterations++
		buf := make([]byte, s.bufferSize)

		read, err := reader.Read(buf)
		if err != nil && read == 0 {
			if err == io.EOF {
				break
			}

			s.log.Errorf("method=%s msg=%s reason=%s pathToFile=%s", "ReadFile", "read file", err.Error(), pathToFile)
			return nil, err
		}

		// Remove unexpected charters, time to time bufio adds many NUL chars at the end of buf
		buf = buf[:read]

		// in case if we stooped on middle of the line because of buffer size, we should read until end of this line
		nextUntillNewline, err := reader.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		words := strings.Split(string(buf), "\n")

		for _, word := range words {
			uniqueWords[word] = struct{}{}
		}
	}

	s.log.Debugf("msg=%s iterations=%d", "finished reading of file", numberOfIterations)

	var res []string
	for word := range uniqueWords {
		res = append(res, word)
	}

	return res, nil
}
