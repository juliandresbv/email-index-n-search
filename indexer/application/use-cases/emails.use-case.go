package usecases

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	customlogger "pkg/custom-logger"
	documentdbclient "pkg/zinc-search-db/db-client/document-db-client"

	"indexer/application/types"
	"indexer/application/utils"
)

var logger = customlogger.NewLogger()

type EmailsUseCase struct {
	documentDbClient documentdbclient.IDocumentDbClient
}

type IEmailsUseCase interface {
	BulkLoadEmails(datasetDataContext types.DatasetDataContext) error
}

func NewEmailsUseCase(documentDbClient documentdbclient.IDocumentDbClient) IEmailsUseCase {
	return &EmailsUseCase{
		documentDbClient,
	}
}

func (emailsUseCase *EmailsUseCase) BulkLoadEmails(datasetDataContext types.DatasetDataContext) error {
	dataDirPath := "./data"
	jsonFilesDirPath := filepath.Join(dataDirPath, datasetDataContext.JsonFilesDirName)

	jsonFilesDirExists := utils.PathExists(jsonFilesDirPath)

	if !jsonFilesDirExists {
		errorStr := "json files for bulk load not found"
		logger.Println(errorStr)

		return errors.New(errorStr)
	}

	bulkLoadResultFilePath := filepath.Join(dataDirPath, datasetDataContext.ZincSearchBulkLoadResultFileName)

	bulkLoadResultFileExists := utils.PathExists(bulkLoadResultFilePath)

	if bulkLoadResultFileExists {
		logger.Println("emails previously bulk loaded into DB")

		return nil
	}

	bulkLoadResultFile, err := os.Create(bulkLoadResultFilePath)

	if err != nil {
		logger.Println("Error creating bulk load result file: ", err)

		return err
	}

	defer bulkLoadResultFile.Close()

	bulkLoadResultFile.WriteString("\"responseStatusCode\",\"responseDataMessage\",\"jsonFile\"\n")

	logger.Println("bulk loading emails into DB...")

	jsonFilesPaths, _ := utils.GetFilesPaths(jsonFilesDirPath, ".json")
	numJsonFiles := len(jsonFilesPaths)
	recordsPerJsonFile := 1000
	numApproxTotalEmails := numJsonFiles * recordsPerJsonFile

	if numJsonFiles <= 0 {
		errorStr := "no json files found for bulk load"
		logger.Println(errorStr)

		return errors.New(errorStr)
	}

	goRoutinesProportion := 0.05
	numGoRoutines := int(float64(numJsonFiles) * goRoutinesProportion)

	var wg sync.WaitGroup
	sem := make(chan struct{}, numGoRoutines)

	millisPerRequest := 750
	throttle := time.NewTicker(time.Duration(millisPerRequest) * time.Millisecond).C

	for index, jsonFilePath := range jsonFilesPaths {
		sem <- struct{}{}
		<-throttle

		wg.Add(1)
		go emailsUseCase.bulkLoadJsonFile(jsonFilePath, numJsonFiles, bulkLoadResultFile, sem, &wg)

		numLoadedFiles := index + 1
		numLoadedEmails := numLoadedFiles * recordsPerJsonFile

		if numLoadedEmails%(numGoRoutines*recordsPerJsonFile) == 0 || numLoadedFiles == numJsonFiles {
			logger.Printf("%d out of ~%d emails loaded into DB\n", numLoadedEmails, numApproxTotalEmails)
		}
	}

	wg.Wait()

	logger.Println("emails bulk loaded into DB successfully!")

	return nil
}

func (emailsUseCase *EmailsUseCase) bulkLoadJsonFile(
	jsonFilePath string,
	numJsonFiles int,
	bulkLoadResultFile *os.File,
	sem chan struct{},
	wg *sync.WaitGroup,
) error {
	defer wg.Done()
	defer func() { <-sem }()

	file, err := os.Open(jsonFilePath)

	if err != nil {
		logger.Println("Error opening file: ", err)

		return err
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		logger.Println("Error reading file: ", err)

		return err
	}

	response, err := emailsUseCase.documentDbClient.BulkV2(data)

	jsonFileName := filepath.Base(jsonFilePath)

	bulkLoadResultFile.WriteString(
		fmt.Sprintf(
			"\"%v\",\"%v\",\"%v\"\n",
			response.Status,
			strings.TrimSpace(fmt.Sprintf("%v %v", response.Data, response.Message)),
			jsonFileName,
		),
	)

	if err != nil {
		logger.Println("Error bulk loading json file: ", err)

		return err
	}

	return nil
}
