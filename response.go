package pinger

// Response is the Response retruned from a ping request
type Response struct {
	StatusCode int
	Body       []byte
	Error      error
}
