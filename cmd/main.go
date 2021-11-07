package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mmtretiak/siden/internal/file/handler"
	"github.com/mmtretiak/siden/internal/file/service"
	"github.com/mmtretiak/siden/pkg/logger"
)

const (
	loggerType  = logger.DevelopmentLogType
	bufferSize  = 1024 * 1024 * 50
	pathToStore = "./store"
	maxMemory   = 1024 * 1024 * 50
	addr        = "localhost:8080"
)

func main() {
	log, err := logger.New(loggerType)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	fileService := service.New(bufferSize, pathToStore, log)
	fileHandler := handler.New(fileService, maxMemory, log)

	router.HandleFunc("/file", fileHandler.Save).Methods(http.MethodPost)
	router.HandleFunc("/file", fileHandler.GetProcessed).Methods(http.MethodGet)

	srv := &http.Server{
		Handler: router,
		Addr:    addr,
	}

	log.Infof("listening on %s", addr)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
