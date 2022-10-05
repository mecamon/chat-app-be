package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, statusCode int, output interface{}) error {
	var i interface{}
	i = output

	if _, ok := i.([]byte); ok {
		return errors.New("output cannot be a []byte, because a conversion to the same type takes place within the function")
	}

	var body []byte
	var err error

	if output != nil {
		body, err = json.Marshal(output)
		if err != nil {
			return err
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)

	return nil
}
