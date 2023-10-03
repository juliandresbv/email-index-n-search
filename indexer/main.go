package main

import (
	documentdbclient "pkg/zinc-search-db/db-client/document-db-client"
	dbclienttypes "pkg/zinc-search-db/db-client/types"

	usecases "indexer/application/use-cases"
	"indexer/application/utils"
	"indexer/framework/config"
	"indexer/framework/profiling"
)

func init() {
	config.LoadEnvVars()
}

func main() {
	defer profiling.SetupProfiling().Stop()

	utils.DownloadAndDecompressDataset()
	utils.DatasetToJsonFiles()

	documentDbClient := documentdbclient.NewDocumentDbClient(dbclienttypes.DbClientConfig{})
	emailsUseCase := usecases.NewEmailsUseCase(documentDbClient)
	emailsUseCase.BulkLoadEmails()
}
