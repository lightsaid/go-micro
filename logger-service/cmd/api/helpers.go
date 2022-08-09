package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

// 定义公共方法

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error{
	maxBytes := 1 << 20 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// 再次解码，验证是否单个json, (防止：{}{} 出现， Decode 每次只解析一个json)
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("读取json错误：只接受一个 json 值")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter,status int, data any, headers ...http.Header) error {
	out,err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, val := range headers[0]{
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil

}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}