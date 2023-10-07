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

	datasetFileName := "enron_mail_20110402.tgz"
	datasetDataContext, _ := utils.GetDatasetDataContext(datasetFileName)

	utils.DownloadAndDecompressDataset(datasetDataContext)
	utils.DatasetToJsonFiles(datasetDataContext)

	documentDbClient := documentdbclient.NewDocumentDbClient(dbclienttypes.DbClientConfig{})
	emailsUseCase := usecases.NewEmailsUseCase(documentDbClient)
	emailsUseCase.BulkLoadEmails(datasetDataContext)
}
