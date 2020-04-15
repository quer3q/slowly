package pkg

type SlowRequest struct {
	Timeout int64 `json:"timeout"`
}

type SlowResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
