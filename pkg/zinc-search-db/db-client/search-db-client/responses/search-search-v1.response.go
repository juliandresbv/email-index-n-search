package responses

type SearchSearchV1Response struct {
	Hits    HitsSearchSearchV1Response   `json:"hits"`
	Shards  ShardsSearchSearchV1Response `json:"_shards"`
	TimeOut bool                         `json:"timed_out"`
	Took    int                          `json:"took"`
	Error   string                       `json:"error,omitempty"`
}

type HitsSearchSearchV1Response struct {
	Hits     []HitsHitsSearchSearchV1Response `json:"hits"`
	MaxScore float64                          `json:"max_score"`
	Total    TotalHitsSearchSearchV1Response  `json:"total"`
}

type ShardsSearchSearchV1Response struct {
	Total     int `json:"total"`
	Sucessful int `json:"sucessful"`
	Skipped   int `json:"skipped"`
	Failed    int `json:"failed"`
}

type HitsHitsSearchSearchV1Response struct {
	Id        string                 `json:"_id"`
	Index     string                 `json:"_index"`
	Score     float64                `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Timestamp int64                  `json:"@timestamp"`
	Type      string                 `json:"_type"`
}

type TotalHitsSearchSearchV1Response struct {
	Value int `json:"value"`
}
