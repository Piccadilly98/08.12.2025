package handlers

import (
	"encoding/json"
	"fmt"
	"time"
)

func GetJsonError(err error) ([]byte, error) {
	str := fmt.Sprintf(`{"status":"error", "text":"%s", "time":"%s"}`, err.Error(), time.Now().String())
	b, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return b, nil
}
