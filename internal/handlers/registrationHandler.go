package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
	linkchecker "github.com/Piccadilly98/linksChecker/internal/linkChecker"
	"github.com/Piccadilly98/linksChecker/internal/storage"
)

type RegistrationHandler struct {
	st *storage.Storage
	lp *linkchecker.LinkProcessor
}

func MakeRegistrationHandler(st *storage.Storage, lp *linkchecker.LinkProcessor) *RegistrationHandler {
	if st == nil || lp == nil {
		return nil
	}
	return &RegistrationHandler{
		st: st,
		lp: lp,
	}
}

func (rh *RegistrationHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	links := rh.readBodyAndValidation(r)
	if links == nil {
		errFmt := fmt.Errorf("invalid body")
		ProcessingError(w, r, errFmt, nil)
		return
	}
	result := rh.lp.LinkChecker(links)
	id := rh.st.RegistrationLinks(result)
	resp := dto.CreateGetInfoBucketDTO(rh.st.GetLiinksInfo(id), id)
	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (rh *RegistrationHandler) readBodyAndValidation(r *http.Request) []string {
	links := &dto.RegistrationLinks{}

	err := json.NewDecoder(r.Body).Decode(links)
	if err != nil {
		return nil
	}
	if !links.Validate() {
		return nil
	}
	links.ProcessingDTO()
	return links.Links
}
