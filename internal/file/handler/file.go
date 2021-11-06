package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mmtretiak/siden/pkg/logger"
)

const (
	formFile      = "file"
	fileNameParam = "fileName"
)

func New(service FileService, maxMemory int64, log logger.Logger) *handler {
	return &handler{
		fileService: service,
		log:         log,
		maxMemory:   maxMemory,
	}
}

type handler struct {
	fileService FileService
	log         logger.Logger
	maxMemory   int64
}

type FileService interface {
	WriteToFile(fileName string, reader io.Reader) error
	ReadFile(fileName string) ([]string, error)
}

func (h *handler) Save(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	err := request.ParseMultipartForm(h.maxMemory)
	if err != nil {
		h.log.Errorf("method=%s msg=%s reason=%s", "Save", "parse multipart form", err.Error())
		http.Error(w, "Failed to parse multipart form.", http.StatusInternalServerError)

		return
	}

	defer request.MultipartForm.RemoveAll()

	file, info, err := request.FormFile(formFile)
	if err != nil {
		h.log.Errorf("method=%s msg=%s reason=%s", "Save", "get file from form", err.Error())
		http.Error(w, "Failed to get file from form.", http.StatusInternalServerError)

		return
	}

	defer file.Close()

	if err := h.fileService.WriteToFile(info.Filename, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetProcessed(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fileName := request.URL.Query().Get(fileNameParam)
	if fileName == "" {
		msg := fmt.Sprintf("missed file name in url query, key: %s", fileNameParam)
		http.Error(w, msg, http.StatusInternalServerError)

		return
	}

	fileContent, err := h.fileService.ReadFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(fileContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		h.log.Errorf("method=%s msg=%s reason=%s", "GetProcessed", "write response", err.Error())
	}
}
