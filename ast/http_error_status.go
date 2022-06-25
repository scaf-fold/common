package ast

// HttpErrorStatus Http Error Status Define
type HttpErrorStatus interface {
	// ServiceCode Need Client self define
	ServiceCode() int
	// ErrMsg will be generated
	ErrMsg() ErrorMsg
	// StatusCode will be generated
	StatusCode() int
}

type ErrorMsg struct {
	Code       int          `json:"code"`
	ServiceErr ServiceError `json:"serviceErr"`
}

type ServiceError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}
