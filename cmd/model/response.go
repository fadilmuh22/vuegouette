package model

type BasicResponse struct {
	Success bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}