package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type CustomErrorMessage interface {
	ErrorMessage() string
}

type MyCustomError struct {
	msg string
}

func (e *MyCustomError) Error() string {
	return e.msg
}

func (e *MyCustomError) ErrorMessage() string {
	return "Custom message: " + e.msg
}

func readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	d := json.NewDecoder(r.Body)
	err := d.Decode(data)
	if err != nil {
		return err
	}

	err = d.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	o, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(o)
	if err != nil {
		return err
	}

	return nil
}

func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	if customMsg, ok := err.(CustomErrorMessage); ok {
		payload.Message = customMsg.ErrorMessage()
	} else {
		payload.Message = err.Error()
	}

	return writeJSON(w, statusCode, payload)
}
