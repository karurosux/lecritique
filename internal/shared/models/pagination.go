package models

// PageRequest represents pagination request parameters
type PageRequest struct {
	Page    int    `json:"page" query:"page"`
	Limit   int    `json:"limit" query:"limit"`
	Sort    string `json:"sort" query:"sort"`
	Order   string `json:"order" query:"order"`
	Filters map[string]interface{} `json:"filters"`
}

// PageResponse represents paginated response
type PageResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

// NewPageRequest creates a new page request with defaults
func NewPageRequest() *PageRequest {
	return &PageRequest{
		Page:  1,
		Limit: 20,
	}
}

// Offset calculates the database offset
func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

// NewPageResponse creates a new page response
func NewPageResponse(data interface{}, page, limit int, total int64) *PageResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	
	return &PageResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}