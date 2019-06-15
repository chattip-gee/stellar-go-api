package http

const Success = "Success"
const Failed = "Failed"

type Status struct {
	Code   int
	Detail string
}

func GetMessage(status *Status) string {
	if status.Code == StatusOK {
		return Success
	} else {
		return status.Detail
	}
}
