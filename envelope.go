package platform

type Envelope struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorEnvelope struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(message string, data interface{}) Envelope {
	return Envelope{Status: CodeSuccess, Message: message, Data: data}
}

func Created(message string, data interface{}) Envelope {
	return Envelope{Status: CodeCreated, Message: message, Data: data}
}

func Fail(code int, message string, errors interface{}) ErrorEnvelope {
	return ErrorEnvelope{Status: code, Message: message, Errors: errors}
}

func BuildPagination(page, size int, total int64) *Pagination {
	if size <= 0 {
		size = 20
	}
	if page <= 0 {
		page = 1
	}
	pages := int(total) / size
	if int(total)%size != 0 {
		pages++
	}
	return &Pagination{Page: page, Size: size, Total: total, TotalPages: pages}
}
