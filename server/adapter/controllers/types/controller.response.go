package types

type SuccessResponse[T interface{}] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
