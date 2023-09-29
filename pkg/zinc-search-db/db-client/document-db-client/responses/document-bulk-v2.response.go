package responses

type DocumentBulkV2Response struct {
	Message     string `json:"message"`
	RecordCount int    `json:"record_count"`
	Error       string `json:"error,omitempty"`
}
