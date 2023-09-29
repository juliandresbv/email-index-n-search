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

	"indexer/application/utils"
)

var logger = customlogger.NewLogger()

type EmailsUseCase struct {
	documentDbClient documentdbclient.IDocumentDbClient
}

type IEmailsUseCase interface {
	BulkLoadEmails() error
}

func NewEmailsUseCase(documentDbClient documentdbclient.IDocumentDbClient) IEmailsUseCase {
	return &EmailsUseCase{
		documentDbClient,
	}
}

func (emailsUseCase *EmailsUseCase) BulkLoadEmails() error {
	jsonFilesPath := "./data/enron_mail_20110402_json"
	jsonFilesDirExists := utils.PathExists(jsonFilesPath)

	if !jsonFilesDirExists {
		errorStr := "json files for bulk load not found"

		logger.Println(errorStr)

		return errors.New(errorStr)
	}

	bulkLoadResultFilePath := "./data/zs-bulkload-result.csv"
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

	bulkLoadResultFile.WriteString("\"responseStatusCode\",\"responseDataMessage\",\"jsonFile\"\n")

	defer bulkLoadResultFile.Close()

	logger.Println("bulk loading emails into DB...")

	jsonFilesPaths, _ := utils.GetFilesPaths(jsonFilesPath, ".json")
	numJsonFiles := len(jsonFilesPaths)

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

		if (numLoadedFiles)%numGoRoutines == 0 || (numLoadedFiles) == numJsonFiles {
			logger.Printf("%d of %d json files loaded into DB\n", numLoadedFiles, numJsonFiles)
		}
	}

	wg.Wait()

	logger.Println("json files bulk loaded into DB")

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
		return err
	}

	<-sem

	return nil
}
