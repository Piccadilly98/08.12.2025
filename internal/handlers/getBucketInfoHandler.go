package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
	document_worker "github.com/Piccadilly98/linksChecker/internal/document_worker"
	"github.com/Piccadilly98/linksChecker/internal/storage"
)

type GetBucketInfoHandler struct {
	st *storage.Storage
}

func MakeGetBucketInfoHandler(st *storage.Storage) *GetBucketInfoHandler {
	return &GetBucketInfoHandler{
		st: st,
	}
}

func (g *GetBucketInfoHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	numbers := g.readBodyAndGetNumbers(r)
	if numbers == nil {
		errFmt := fmt.Errorf("invalid body")
		ProcessingError(w, r, errFmt, nil, http.StatusBadRequest)
		return
	}
	res, err := g.st.GetBucketsInfo(numbers...)
	if err != nil {
		ProcessingError(w, r, err, nil, http.StatusBadRequest)
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

func (g *GetBucketInfoHandler) readBodyAndGetNumbers(r *http.Request) []int64 {
	dto := &dto.InfoWithNumbersBucketDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		return nil
	}
	if !dto.Validate() {
		return nil
	}
	dto.ProcessingDTO()
	return dto.LinksList
}
