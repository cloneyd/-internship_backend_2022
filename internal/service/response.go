package service

type errorResponse struct {
	Error string `json:"error"`
}

type resultResponse struct {
	Result string `json:"result"`
}

func newResultResponse(result string) *resultResponse {
	return &resultResponse{
		Result: result,
	}
}

func newErrorResponse(err error) *errorResponse {
	return &errorResponse{
		Error: err.Error(),
	}
}
