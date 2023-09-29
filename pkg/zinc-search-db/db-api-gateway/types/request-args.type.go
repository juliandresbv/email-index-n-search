package types

type RequestArgs struct {
	BytesBody []byte
	Body      map[string]interface{}
	Params    []string
	Query     map[string]interface{}
}
