package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	customlogger "pkg/custom-logger"

	"indexer/application/types"
)

var loggerDatasetUtil = customlogger.NewLogger()

var emailMetadataHeadersKeys = []string{
	"MessageId",
	"Date",
	"From",
	"To",
	"Subject",
	"Cc",
	"MimeVersion",
	"ContentType",
	"ContentTransferEncoding",
	"Bcc",
	"XFrom",
	"XTo",
	"XCc",
	"XBcc",
	"XFolder",
	"XOrigin",
	"XFileName",
}

var emailMetadataHeadersMap = map[string]string{
	"MessageId":               "Message-ID:",
	"Date":                    "Date:",
	"From":                    "From:",
	"To":                      "To:",
	"Subject":                 "Subject:",
	"Cc":                      "Cc:",
	"MimeVersion":             "Mime-Version:",
	"ContentType":             "Content-Type:",
	"ContentTransferEncoding": "Content-Transfer-Encoding:",
	"Bcc":                     "Bcc:",
	"XFrom":                   "X-From:",
	"XTo":                     "X-To:",
	"XCc":                     "X-cc:",
	"XBcc":                    "X-bcc:",
	"XFolder":                 "X-Folder:",
	"XOrigin":                 "X-Origin:",
	"XFileName":               "X-FileName:",
}

var ignoredMetadataHeadersKeys = []string{
	"MimeVersion",
	"ContentType",
	"ContentTransferEncoding",
}

func DownloadAndDecompressDataset(datasetDataContext types.DatasetDataContext) {
	dataDirPath := "./data"

	isDataDirCreated := PathExists(dataDirPath)

	if !isDataDirCreated {
		loggerDatasetUtil.Println("creating data directory...")

		err := os.MkdirAll(dataDirPath, 0755)

		if err != nil {
			loggerDatasetUtil.Println("error creating data directory: ", err)

			panic(err)
		}

		loggerDatasetUtil.Println("data directory created successfully!")
	} else {
		loggerDatasetUtil.Println("data directory previously created")
	}

	datasetFilePath := filepath.Join(dataDirPath, datasetDataContext.DatasetFileName)

	isDatasetDownloaded := PathExists(datasetFilePath)

	if !isDatasetDownloaded {
		loggerDatasetUtil.Println("downloading dataset...")

		datasetUrl := "https://www.cs.cmu.edu/~./enron/" + datasetDataContext.DatasetFileName

		err := downloadFile(datasetUrl, datasetDataContext)

		if err != nil {
			loggerDatasetUtil.Println("error downloading dataset: ", err)

			panic(err)
		}

		loggerDatasetUtil.Println("dataset downloaded successfully!")
	} else {
		loggerDatasetUtil.Println("dataset previously downloaded")
	}

	decompressedDatasetDirPath := filepath.Join(dataDirPath, datasetDataContext.DatasetFileNameNoExt)

	isDatasetDecompressed := PathExists(decompressedDatasetDirPath)

	if !isDatasetDecompressed {
		loggerDatasetUtil.Println("decompressing dataset...")

		err := decompressTgzFile(datasetDataContext)

		if err != nil {
			loggerDatasetUtil.Println("error decompressing dataset: ", err)

			panic(err)
		}

		loggerDatasetUtil.Println("dataset decompressed successfully!")
	} else {
		loggerDatasetUtil.Println("dataset previously decompressed")
	}
}

func DatasetToJsonFiles(datasetDataContext types.DatasetDataContext) {
	dataDirPath := "./data"
	jsonFilesDirPath := filepath.Join(dataDirPath, datasetDataContext.JsonFilesDirName)

	jsonFilesDirExists := PathExists(jsonFilesDirPath)

	if !jsonFilesDirExists {
		err := os.MkdirAll(jsonFilesDirPath, 0755)

		if err != nil {
			loggerDatasetUtil.Println("error creating json files directory: ", err)

			panic(err)
		}
	}

	jsonFilesPaths, err := GetFilesPaths(jsonFilesDirPath, ".json")

	if err != nil {
		loggerDatasetUtil.Println("error getting json files paths: ", err)

		panic(err)
	}

	numExistingJsonFiles := len(jsonFilesPaths)

	if numExistingJsonFiles > 0 {
		loggerDatasetUtil.Println("json files previously created")

		return
	}

	decompressedDatasetDirPath := filepath.Join(dataDirPath, datasetDataContext.DatasetFileNameNoExt)

	datasetFilesPaths, err := GetFilesPaths(decompressedDatasetDirPath, ".")

	if err != nil {
		loggerDatasetUtil.Println("error getting dataset files paths: ", err)

		panic(err)
	}

	loggerDatasetUtil.Println("creating json files...")

	chunkSize := 1000
	numFiles := len(datasetFilesPaths)
	numJsonFilesToCreate := math.Ceil(float64(numFiles) / float64(chunkSize))

	goRoutinesProportion := 0.05
	numGoRoutines := int(numJsonFilesToCreate * goRoutinesProportion)

	var wg sync.WaitGroup
	sem := make(chan struct{}, numGoRoutines)

	for i := 0; i < numFiles; i += chunkSize {
		sem <- struct{}{}

		chunkEnd := int(math.Floor(
			math.Min(
				float64(i+chunkSize),
				float64(numFiles),
			),
		))

		filesPathsChunk := datasetFilesPaths[i:chunkEnd]
		chunkId := (i / chunkSize) + 1

		wg.Add(1)
		go datasetChunkToJsonFile(filesPathsChunk, chunkId, jsonFilesDirPath, sem, &wg)

		if chunkEnd%(numGoRoutines*chunkSize) == 0 || chunkEnd == numFiles {
			loggerDatasetUtil.Printf("%v out of %v emails to json file\n", chunkEnd, numFiles)
		}
	}

	wg.Wait()

	loggerDatasetUtil.Println("json files created successfully!")
}

func GetDatasetDataContext(datasetFileName string) types.DatasetDataContext {
	datasetFileNameHasExt := strings.Contains(datasetFileName, ".")

	if !datasetFileNameHasExt {
		errorStr := "dataset file name must have an extension"
		loggerDatasetUtil.Println(errorStr)

		panic(errors.New(errorStr))
	}

	datasetFileNameNoExt := strings.Split(datasetFileName, ".")[0]
	jsonFilesDirName := datasetFileNameNoExt + "_json"
	zincSearchBulkLoadResultFileName := datasetFileNameNoExt + "_zs_bulkload_result.csv"

	return types.DatasetDataContext{
		DatasetFileName:                  datasetFileName,
		DatasetFileNameNoExt:             datasetFileNameNoExt,
		JsonFilesDirName:                 jsonFilesDirName,
		ZincSearchBulkLoadResultFileName: zincSearchBulkLoadResultFileName,
	}
}

func downloadFile(
	datasetUrl string,
	datasetDataContext types.DatasetDataContext,
) error {
	resp, err := http.Head(datasetUrl)

	if err != nil {
		loggerDatasetUtil.Println("error getting dataset metadata from URL: ", err)

		return err
	}

	defer resp.Body.Close()

	contentSizeBytes := int(resp.ContentLength)

	oneMegaByte := 1024 * 1024
	chunkSizeBytes := 5 * oneMegaByte

	numChunks := int(math.Ceil(float64(contentSizeBytes) / float64(chunkSizeBytes)))

	chunkSize := contentSizeBytes / numChunks

	datasetFilePath := filepath.Join("./data", datasetDataContext.DatasetFileName)
	file, err := os.Create(datasetFilePath)

	if err != nil {
		loggerDatasetUtil.Println("error creating file to download: ", err)

		return err
	}

	defer file.Close()

	var wg sync.WaitGroup

	for i := 0; i < numChunks; i++ {
		startByte := i * chunkSize
		endByte := startByte + chunkSize

		if i == numChunks-1 {
			endByte = contentSizeBytes
		}

		wg.Add(1)
		go downloadFileChunk(datasetUrl, startByte, endByte, file, &wg)
	}

	wg.Wait()

	return nil
}

func downloadFileChunk(
	url string,
	startByte int,
	endByte int,
	writer io.WriterAt,
	wg *sync.WaitGroup,
) error {
	defer wg.Done()

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rangeHeader := fmt.Sprintf("bytes=%d-%d", startByte, endByte)
	req.Header.Add("Range", rangeHeader)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		loggerDatasetUtil.Println("error downloading chunk: ", err)

		return err
	}

	defer resp.Body.Close()

	chunkSizeBytes := 6 * (1024 * 1024)
	buf := make([]byte, chunkSizeBytes)

	for {
		n, err := resp.Body.Read(buf)

		if err != nil && err != io.EOF {
			loggerDatasetUtil.Println("error reading file to download: ", err)

			return err
		}
		if n == 0 {
			break
		}

		_, err = writer.WriteAt(buf[:n], int64(startByte))

		if err != nil {
			loggerDatasetUtil.Println("error writing file to download: ", err)

			return err
		}

		startByte += n
	}

	return nil
}

func decompressTgzFile(datasetDataContext types.DatasetDataContext) error {
	dataDirPath := "./data"
	datasetFilePath := filepath.Join(dataDirPath, datasetDataContext.DatasetFileName)

	tgzFile, err := os.Open(datasetFilePath)

	if err != nil {
		loggerDatasetUtil.Println("error opening tgz file: ", err)

		return err
	}

	defer tgzFile.Close()

	gzReader, err := gzip.NewReader(tgzFile)

	if err != nil {
		loggerDatasetUtil.Println("error creating gzip reader: ", err)

		return err
	}

	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		tarHeader, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case tarHeader == nil:
			continue
		}

		var target string
		compressedFileOriginalFirstDir := strings.Split(tarHeader.Name, "/")[0]

		if compressedFileOriginalFirstDir != datasetDataContext.DatasetFileNameNoExt {
			target = filepath.Join(dataDirPath, datasetDataContext.DatasetFileNameNoExt, tarHeader.Name)
		} else {
			target = filepath.Join(dataDirPath, tarHeader.Name)
		}

		switch tarHeader.Typeflag {
		case tar.TypeDir:
			_, err := os.Stat(target)

			if err != nil {
				err := os.MkdirAll(target, os.FileMode(tarHeader.Mode))

				if err != nil {
					loggerDatasetUtil.Println("error creating directory: ", err)

					return err
				}
			}
		case tar.TypeReg:
			err := os.MkdirAll(filepath.Dir(target), 0755)

			if err != nil {
				loggerDatasetUtil.Println("error creating directory: ", err)

				return err
			}

			file, err := os.Create(target)

			if err != nil {
				loggerDatasetUtil.Println("error creating file: ", err)

				return err
			}

			_, err = io.Copy(file, tarReader)

			if err != nil {
				loggerDatasetUtil.Println("error copying file: ", err)

				file.Close()

				return err
			}

			file.Close()
		}
	}
}

func datasetChunkToJsonFile(
	filesPathsChunk []string,
	chunkId int,
	jsonFilesDirPath string,
	sem chan struct{},
	wg *sync.WaitGroup,
) error {
	defer wg.Done()
	defer func() { <-sem }()

	err := writeJsonFile(filesPathsChunk, chunkId, jsonFilesDirPath)

	if err != nil {
		loggerDatasetUtil.Println("error writing json file: ", err)

		return err
	}

	return nil
}

func writeJsonFile(
	filesPathsChunk []string,
	chunkId int,
	jsonFilesDirPath string,
) error {
	jsonFileName := fmt.Sprintf("emails-%v.json", chunkId)

	jsonFile, err := os.Create(filepath.Join(jsonFilesDirPath, jsonFileName))

	if err != nil {
		loggerDatasetUtil.Println("error creating json file: ", err)

		return err
	}

	defer jsonFile.Close()

	var fileData *os.File

	defer fileData.Close()

	numFilesChunk := len(filesPathsChunk)
	indexName := "emails"

	jsonFile.WriteString("{\n")
	jsonFile.WriteString(fmt.Sprintf("\"index\": \"%v\",\n", indexName))
	jsonFile.WriteString("\"records\": [\n")

	for index, path := range filesPathsChunk {
		fileData, err = os.Open(path)

		if err != nil {
			loggerDatasetUtil.Println("error opening file: ", err)

			return err
		}

		fileDataScanner := bufio.NewScanner(fileData)

		emailMap := map[string]string{}

		hasBodyStarted := false
		currMetadataKey := ""

		for fileDataScanner.Scan() {
			line := fileDataScanner.Text()

			isMetadataLine := false

			if !hasBodyStarted {
				for _, key := range emailMetadataHeadersKeys {
					lineToCheckMetadataHeaderPrefix := line
					lineToCheckMetadataHeaderPrefix = strings.TrimSpace(lineToCheckMetadataHeaderPrefix)

					if currMetadataKey == emailMetadataHeadersKeys[len(emailMetadataHeadersKeys)-1] &&
						len(lineToCheckMetadataHeaderPrefix) <= 0 {
						hasBodyStarted = true

						break
					}

					emailMetadataHeadersValue := emailMetadataHeadersMap[key]

					if strings.HasPrefix(lineToCheckMetadataHeaderPrefix, emailMetadataHeadersValue) {
						isMetadataLine = true
						currMetadataKey = key

						break
					}
				}
			}

			if !hasBodyStarted {
				if slices.Contains(ignoredMetadataHeadersKeys, currMetadataKey) {
					continue
				} else {
					lineContent := line

					if isMetadataLine {
						lineMetadataHeader := emailMetadataHeadersMap[currMetadataKey]

						lineContent = strings.TrimSpace(lineContent)
						lineContent = lineContent[len(lineMetadataHeader):]
						lineContent = strings.TrimSpace(lineContent)

						if len(lineContent) <= 0 {
							continue
						}

						emailMap[currMetadataKey] = lineContent
					} else {
						lineContent = strings.TrimSpace(lineContent)

						if len(lineContent) <= 0 {
							continue
						}

						emailMap[currMetadataKey] += lineContent
					}
				}
			} else {
				emailMap["Body"] += line
			}
		}

		emailMapJson, err := json.Marshal(emailMap)

		if err != nil {
			loggerDatasetUtil.Println("error marshaling email map: ", err)

			return err
		}

		if index < numFilesChunk-1 {
			jsonFile.WriteString(fmt.Sprintf("%v,\n", string(emailMapJson)))
		} else {
			jsonFile.WriteString(fmt.Sprintf("%v\n", string(emailMapJson)))
		}
	}

	jsonFile.WriteString("]\n")
	jsonFile.WriteString("}\n")

	return nil
}
