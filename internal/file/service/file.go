package service

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmtretiak/siden/pkg/logger"
)

func New(bufferSize int, pathToStore string, log logger.Logger) *service {
	return &service{
		bufferSize:  bufferSize,
		log:         log,
		pathToStore: pathToStore,
	}
}

type service struct {
	bufferSize  int
	log         logger.Logger
	pathToStore string
}

func (s *service) WriteToFile(fileName string, reader io.Reader) error {
	file, err := os.OpenFile(filepath.Join(s.pathToStore, fileName), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s fileName=%s", "WriteToFile", "open file", err.Error(), fileName)
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	buf := make([]byte, s.bufferSize)

	for {
		// read a chunk
		read, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			s.log.Errorf("method=%s msg=%s reason=%s", "WriteToFile", "read chunk", err.Error())
			return err
		}

		if read == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:read]); err != nil {
			s.log.Errorf("method=%s msg=%s reason=%s", "WriteToFile", "write chunk", err.Error())
			return err
		}
	}

	if err = w.Flush(); err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s", "WriteToFile", "flush writer", err.Error())
		return err
	}

	return nil
}

func (s *service) ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(filepath.Join(s.pathToStore, fileName))
	if err != nil {
		s.log.Errorf("method=%s msg=%s reason=%s fileName=%s", "ReadFile", "open file", err.Error(), fileName)
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	uniqueWords := make(map[string]struct{})

	numberOfIterations := 0
	buf := make([]byte, s.bufferSize)

	for {
		read, err := reader.Read(buf)
		if err != nil && read == 0 {
			if err == io.EOF {
				break
			}

			s.log.Errorf("method=%s msg=%s reason=%s fileName=%s", "ReadFile", "read file", err.Error(), fileName)
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

		numberOfIterations++
	}

	s.log.Debugf("msg=%s iterations=%d", "finished reading of file", numberOfIterations)

	var res []string
	for word := range uniqueWords {
		// TODO find why empty string presented in map
		if word == "" {
			continue
		}

		res = append(res, word)
	}

	return res, nil
}
