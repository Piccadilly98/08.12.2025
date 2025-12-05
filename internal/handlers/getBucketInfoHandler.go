package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
	document_worker "github.com/Piccadilly98/linksChecker/internal/documentWorker.go"
	"github.com/Piccadilly98/linksChecker/internal/storage"
)

type getBucketInfoHandler struct {
	st *storage.Storage
}

func MakeGetBucketInfoHandler(st *storage.Storage) *getBucketInfoHandler {
	return &getBucketInfoHandler{st: st}
}

func (g *getBucketInfoHandler) Hadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	numbers := g.readBodyAndGetNumbers(r)
	if numbers == nil {
		errFmt := fmt.Errorf("invalid body")
		b := dto.MakeResponseDTO(errFmt, nil)
		resp, err := json.Marshal(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	res, err := g.st.GetBucketsInfo(numbers...)
	if err != nil {
		b := dto.MakeResponseDTO(err, nil)
		resp, err := json.Marshal(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	b, err := document_worker.CreateDocument(res)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(b)
}

func (g *getBucketInfoHandler) readBodyAndGetNumbers(r *http.Request) []int64 {
	dto := &dto.InfoWithNumbersBucketDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		return nil
	}
	if !dto.Validate() {
		return nil
	}
	return dto.LinksList
}
