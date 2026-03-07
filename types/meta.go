package types

type Meta struct {
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
	Total     int64  `json:"total,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}
