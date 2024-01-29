package errs

type ErrorResponse struct {
	Error map[string]string `json:"error"`
}
