package responder

type Pagination struct {
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	Sort       string      `json:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Keyword    string      `json:"keyword"`
	Rows       interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
