package documentdbclient

import (
	"encoding/json"
	"net/http"

	dbapigateway "pkg/zinc-search-db/db-api-gateway"
	dbapigatewaytypes "pkg/zinc-search-db/db-api-gateway/types"
	"pkg/zinc-search-db/db-client/document-db-client/responses"
	dbclienttypes "pkg/zinc-search-db/db-client/types"
)

type DocumentDbClient struct {
	dbApiGateway      dbapigateway.IDbApiGateway
	endpointResources map[string]string
}

type IDocumentDbClient interface {
	BulkV2(
		documentBulkV2Dto []byte,
	) (dbapigatewaytypes.ApiResponse[responses.DocumentBulkV2Response], error)
}

func NewDocumentDbClient(dbClientConfig dbclienttypes.DbClientConfig) IDocumentDbClient {
	dbApiGateway := dbClientConfig.DbApiGateway

	if dbApiGateway == nil {
		dbApiGateway = dbapigateway.NewDbApiGateway()
	}

	endpointResources := map[string]string{
		"BulkV2": "_bulkv2",
	}

	return &DocumentDbClient{
		dbApiGateway,
		endpointResources,
	}
}

func (documentClient *DocumentDbClient) BulkV2(
	documentBulkV2Dto []byte,
) (dbapigatewaytypes.ApiResponse[responses.DocumentBulkV2Response], error) {
	response, err := documentClient.dbApiGateway.MakeDbApiRequest(
		http.MethodPost,
		documentClient.endpointResources["BulkV2"],
		dbapigatewaytypes.RequestArgs{
			BytesBody: documentBulkV2Dto,
		},
	)

	var documentBulkV2Response responses.DocumentBulkV2Response
	marshaledResponseData, _ := json.Marshal(response.Data)
	json.Unmarshal(marshaledResponseData, &documentBulkV2Response)

	typedResponse := dbapigatewaytypes.ApiResponse[responses.DocumentBulkV2Response]{
		Status:  response.Status,
		Message: response.Message,
		Data:    documentBulkV2Response,
	}

	return typedResponse, err
}
