package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func decodeJSON(
	r *http.Request,
	dst any,
) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err != nil {
			return errors.New("request body must contain a single JSON object")
		}

		return err
	}

	return nil
}

func writeJSONDecodeError(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError

	if errors.As(err, &maxBytesError) {
		http.Error(
			w,
			"request body too large",
			http.StatusRequestEntityTooLarge,
		)
		return
	}

	http.Error(
		w,
		"invalid request body",
		http.StatusBadRequest,
	)
}
