package upload

type HostResponse struct {
	Host       string
	StatusCode int
	Body       map[string]interface{}
	Err        error
}
