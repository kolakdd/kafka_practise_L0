package apiError

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Error backendError `json:"error"`
}

type backendError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type BackendErrorInternal struct {
	Code     int    `json:"code"`
	Text     string `json:"text"`
	HTTPCode int
}

func (e BackendErrorInternal) Error() string {
	return fmt.Sprintf("code: %v, text: %s, httpCode: %v\n", e.Code, e.Text, e.HTTPCode)
}

var (
	BadRequest       = &BackendErrorInternal{0, "Неверная форма запроса", http.StatusBadRequest}
	NotFound         = &BackendErrorInternal{33, "Не найдено", http.StatusNotFound}
	MethodNotAllowed = &BackendErrorInternal{44, "Метод запрещен", http.StatusMethodNotAllowed}
	InternalError    = &BackendErrorInternal{55, "Внутренняя ошибка", http.StatusInternalServerError}
	MarshalError     = &BackendErrorInternal{104, "Ошибка сериализации", http.StatusInternalServerError}
)

func BackendErrorWrite(w http.ResponseWriter, err *BackendErrorInternal) {
	w.WriteHeader(err.HTTPCode)
	jData, _ := json.Marshal(errorResponse{Error: backendError{Code: err.Code, Text: err.Text}})
	_, _ = w.Write(jData)
}
