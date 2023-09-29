package dtos

type DocumentBulkV2Dto struct {
	Index   string        `json:"index"`
	Records []interface{} `json:"records"`
}
