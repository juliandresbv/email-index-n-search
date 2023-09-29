package dbapigateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	customlogger "pkg/custom-logger"
	"pkg/zinc-search-db/db-api-gateway/types"
)

var logger = customlogger.NewLogger()

type DbApiGateway struct {
	baseDbApiRequest *http.Request
}

type IDbApiGateway interface {
	MakeDbApiRequest(
		endpointResource string,
		httpMethod string,
		reqArgs types.RequestArgs,
	) (types.ApiResponse[interface{}], error)
}

func NewDbApiGateway() IDbApiGateway {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" {
		errorStr := "invalid database credentials"

		logger.Println(errorStr)
		panic(errorStr)
	}

	baseDbApiUrl := fmt.Sprintf("http://%v:%v/api", dbHost, dbPort)

	baseDbApiRequest, _ := http.NewRequest(http.MethodHead, baseDbApiUrl, nil)
	baseDbApiRequest.SetBasicAuth(dbUser, dbPassword)
	baseDbApiRequest.Header.Set("Content-Type", "application/json")

	return &DbApiGateway{
		baseDbApiRequest,
	}
}

func (dbApiGateway *DbApiGateway) MakeDbApiRequest(
	httpMethod string,
	endpointResource string,
	args types.RequestArgs,
) (types.ApiResponse[interface{}], error) {
	apiResponse := types.ApiResponse[interface{}]{
		Status: http.StatusInternalServerError,
	}

	apiRequest := dbApiGateway.baseDbApiRequest.Clone(context.TODO())

	apiRequest.Method = httpMethod

	if endpointResource != "" {
		apiRequest.URL.Path = strings.Join([]string{apiRequest.URL.Path, endpointResource}, "/")
	}

	if args.BytesBody != nil {
		apiRequest.Body = io.NopCloser(bytes.NewReader(args.BytesBody))
	} else if args.Body != nil {
		payload, _ := json.Marshal(args.Body)
		apiRequest.Body = io.NopCloser(bytes.NewReader(payload))
	}

	if args.Query != nil {
		query := apiRequest.URL.Query()

		for key, value := range args.Query {
			query.Add(key, fmt.Sprintf("%v", value))
		}

		apiRequest.URL.RawQuery = query.Encode()
	}

	response, err := http.DefaultClient.Do(apiRequest)

	if err != nil {
		logger.Println(err)

		apiResponse.Message = err.Error()

		return apiResponse, err
	}

	defer response.Body.Close()

	apiResponse.Status = response.StatusCode

	resBody, _ := io.ReadAll(response.Body)
	var unmarshaledResBody interface{}
	json.Unmarshal(resBody, &unmarshaledResBody)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorStr := fmt.Sprintf("%v", unmarshaledResBody)
		apiResponse.Message = errorStr

		logger.Println(errorStr)

		return apiResponse, fmt.Errorf(errorStr)
	}

	apiResponse.Data = unmarshaledResBody

	return apiResponse, nil
}
