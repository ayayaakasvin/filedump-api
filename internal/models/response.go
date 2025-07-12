package models

import (
	"encoding/json"
	"net/http"
)

const (
	ContentType = "Content-Type"
	AppJson 	= "application/json" 
)

type JsonResponse struct {
	Status
	Data 	`json:"data,omitempty"`
}

func SendSuccessJson(w http.ResponseWriter, statusCode int, resp Data) error {
	w.Header().Set(ContentType, AppJson)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(NewJsonResponse(StatusSuccess(), resp))
}

func SendErrorJson(w http.ResponseWriter, statusCode int, msg string, args... any) error {
	w.Header().Set(ContentType, AppJson)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(NewJsonResponse(StatusError(msg, args...), nil))
}

func NewJsonResponse(status Status, data Data) *JsonResponse {
	resp := new(JsonResponse)
	resp.Status = status
	resp.Data = data

	return resp
}