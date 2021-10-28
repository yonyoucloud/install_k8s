package handler

type (
	Response struct {
		Code int         `json:"Code"`
		Msg  string      `json:"Msg"`
		Data interface{} `json:"Data"`
	}
)
