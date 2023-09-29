package dtos

type QuerySearchSearchV1Dto struct {
	Term      string `json:"term"`
	Field     string `json:"field"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

type SearchSearchV1Dto struct {
	SearchType string                 `json:"search_type"`
	Query      QuerySearchSearchV1Dto `json:"query"`
	SortFields []string               `json:"sort_fields"`
	From       int                    `json:"from"`
	MaxResults int                    `json:"max_results" binding:"gte=1,lte=100"`
	Source     []string               `json:"_source"`
}
