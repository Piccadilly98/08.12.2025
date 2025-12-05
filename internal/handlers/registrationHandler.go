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

type registrationHandler struct {
	st *storage.Storage
	lp *linkchecker.LinkProcessor
}

func MakeRegistrationHandler(st *storage.Storage, lp *linkchecker.LinkProcessor) *registrationHandler {
	if st == nil || lp == nil {
		return nil
	}
	return &registrationHandler{
		st: st,
		lp: lp,
	}
}

func (rh *registrationHandler) Hadler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	links := rh.readBodyAndValidation(r)
	if links == nil {
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

func (rh *registrationHandler) readBodyAndValidation(r *http.Request) []string {
	links := &dto.RegistrationLinks{}

	err := json.NewDecoder(r.Body).Decode(links)
	if err != nil {
		return nil
	}
	if !links.Validate() {
		return nil
	}
	return links.Links
}
