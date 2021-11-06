package handler

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mock_handler "github.com/mmtretiak/siden/internal/file/handler/mocks"
	"github.com/mmtretiak/siden/pkg/logger"
)

const (
	fileName    = "file_example.txt"
	pathToStore = "../test_data"
	maxMemory   = 0
)

func TestHandler_Save(t *testing.T) {
	ctrl := gomock.NewController(t)

	log, err := logger.New(logger.DevelopmentLogType)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(pathToStore, fileName))
	if err != nil {
		t.Fatal(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	go func() {
		defer writer.Close()
		formFile, err := writer.CreateFormFile(formFile, fileName)
		if err != nil {
			t.Error(err)
		}

		written, err := formFile.Write(fileContent)
		if err != nil {
			t.Error(err)
		}

		if written != len(fileContent) {
			t.Errorf("unexpected bytes written, expected %d, actual %d", len(fileContent), written)
		}
	}()

	req, err := http.NewRequest(http.MethodPost, "/file", pipeReader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	fileService := mock_handler.NewMockFileService(ctrl)
	fileService.EXPECT().WriteToFile(fileName, fileMatcher{expectedContent: fileContent})

	fileHandler := New(fileService, maxMemory, log)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler.Save)

	handler.ServeHTTP(rr, req)

	ctrl.Finish()
}

type fileMatcher struct {
	expectedContent []byte
}

func (e fileMatcher) Matches(x interface{}) bool {
	reader, ok := x.(io.Reader)
	if !ok {
		return false
	}

	fileContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(e.expectedContent, fileContent)
}

func (e fileMatcher) String() string {
	return "matches file content"
}
