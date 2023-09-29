package searchdbclient

import (
	"encoding/json"
	"net/http"
	"strings"

	dbapigateway "pkg/zinc-search-db/db-api-gateway"
	dbapigatewaytypes "pkg/zinc-search-db/db-api-gateway/types"
	"pkg/zinc-search-db/db-client/search-db-client/dtos"
	"pkg/zinc-search-db/db-client/search-db-client/responses"
	dbclienttypes "pkg/zinc-search-db/db-client/types"
)

type SearchDbClient struct {
	dbApiGateway      dbapigateway.IDbApiGateway
	endpointResources map[string]string
}

type ISearchDbClient interface {
	SearchV1(
		indexName string,
		searchSearchV1Dto dtos.SearchSearchV1Dto,
	) (dbapigatewaytypes.ApiResponse[responses.SearchSearchV1Response], error)
}

func NewSearchDbClient(dbClientConfig dbclienttypes.DbClientConfig) ISearchDbClient {
	dbApiGateway := dbClientConfig.DbApiGateway

	if dbApiGateway == nil {
		dbApiGateway = dbapigateway.NewDbApiGateway()
	}

	endpointResources := map[string]string{
		"SearchV1": "_search",
	}

	return &SearchDbClient{
		dbApiGateway,
		endpointResources,
	}
}

func (searchDbClient *SearchDbClient) SearchV1(
	indexName string,
	searchSearchV1Dto dtos.SearchSearchV1Dto,
) (dbapigatewaytypes.ApiResponse[responses.SearchSearchV1Response], error) {
	var body map[string]interface{}
	marshaledBody, _ := json.Marshal(searchSearchV1Dto)
	json.Unmarshal(marshaledBody, &body)

	endpointResource := strings.Join([]string{
		indexName,
		searchDbClient.endpointResources["SearchV1"],
	}, "/")

	response, err := searchDbClient.dbApiGateway.MakeDbApiRequest(
		http.MethodPost,
		endpointResource,
		dbapigatewaytypes.RequestArgs{
			Body: body,
		},
	)

	var searchSearchV1Response responses.SearchSearchV1Response
	marshaledResponseData, _ := json.Marshal(response.Data)
	json.Unmarshal(marshaledResponseData, &searchSearchV1Response)

	typedResponse := dbapigatewaytypes.ApiResponse[responses.SearchSearchV1Response]{
		Status:  response.Status,
		Message: response.Message,
		Data:    searchSearchV1Response,
	}

	return typedResponse, err
}
