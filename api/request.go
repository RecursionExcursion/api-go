package api

import (
	"encoding/json"
	"net/http"
)

func DecodeJson[T any](r *http.Request) (T, error) {
	defer r.Body.Close()

	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	return t, err
}
