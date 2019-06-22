package http

import (
	. "github.com/chattip-gee/stellar-go-api/constant"
)

type Status struct {
	Code   int
	Detail string
}

func GetMessage(status *Status) string {
	if status.Code == StatusOK {
		return SUCCESS
	} else {
		return status.Detail
	}
}
