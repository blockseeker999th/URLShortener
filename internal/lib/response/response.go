package response

type ResponseStruct struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func ResponseWithoutPayload(status int, err string) ResponseStruct {
	return ResponseStruct{
		Status: status,
		Error:  err,
	}
}
