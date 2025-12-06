package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
)

func GetJsonError(err error) ([]byte, error) {
	str := fmt.Sprintf(`{"status":"error", "text":"%s", "time":"%s"}`, err.Error(), time.Now().String())
	b, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ProcessingError(w http.ResponseWriter, r *http.Request, err error, data *string) {
	dto := dto.MakeResponseDTO(err, data)
	b, err := json.Marshal(dto)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(b)
}
