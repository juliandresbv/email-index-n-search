package types

import (
	dbapigateway "pkg/zinc-search-db/db-api-gateway"
)

type DbClientConfig struct {
	DbApiGateway dbapigateway.IDbApiGateway
}
