package types

type Meta struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	Total     int64  `json:"total"`
	RequestID string `json:"request_id"`
}

func BuildMeta(page, perPage, total int, reqID string) *Meta {
	return &Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     int64(total),
		RequestID: reqID,
	}
}
