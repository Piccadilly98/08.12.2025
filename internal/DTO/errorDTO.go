package dto

import "time"

type ResponseDTO struct {
	Status string  `json:"status"`
	Error  *string `json:"error,omitempty"`
	Date   string  `json:"date"`
	Data   *string `json:"data,omitempty"`
}

func MakeResponseDTO(err error, data *string) *ResponseDTO {
	resp := &ResponseDTO{}
	resp.Date = time.Now().String()
	if err != nil {
		resp.Status = "error"
		str := err.Error()
		resp.Error = &str
	} else {
		resp.Status = "ok"
		resp.Data = data
	}

	return resp
}
