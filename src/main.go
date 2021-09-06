package main

import (
	"encoding/json"
	"fmt"
	"github.com/uly55e5/DataciteDoiDownloader/datacite"
	"github.com/uly55e5/DataciteDoiDownloader/graphql"
	"os"
	"time"
)

const filename = "result.json"
const query = "NOT chemistry"
const firstCount = 200
const dataFormat = "DoiData"
const logName = "datacitetool.log"

var logFile *os.File

func main() {
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logError("Could not open log file", err)
	}
	result := datacite.NamedResponse[dataFormat]
	resultNodeLen := 0
	var currentCursor = ""

	err = readResultFile(filename, &result)
	if err == nil {
		if result.ToolData.DataFormat == dataFormat && result.ToolData.Query == query {
			currentCursor = result.PageInfo().EndCursor
			resultNodeLen = len(result.Nodes())
		} else {
			result = datacite.NamedResponse[dataFormat]
		}
	}
	resultNodeStartLen := resultNodeLen
	result.ToolData = datacite.ToolData{Query: query, DataFormat: dataFormat, Time: time.Now()}
	firstStart := time.Now()
	reqData := datacite.NamedRequests[dataFormat]
	reqData.Variables["query"] = query
	reqData.Variables["first"] = firstCount
	for {
		reqStartTime := time.Now()
		response := datacite.NamedResponse[dataFormat]
		reqData.Variables["cursorId"] = currentCursor
		if err := graphql.GetDataCiteGraphQLResult(&reqData, &response); err != nil {
			logError("Could not get GraphQL response", err)
			continue
		}
		reqTime := time.Since(reqStartTime)
		if !response.PageInfo().HasNextPage || response.PageInfo().EndCursor == "" {
			break
		}
		currentCursor = response.PageInfo().EndCursor
		if resultNodeLen > 0 {
			result.Data.Datasets.Nodes = append(result.Nodes(), response.Nodes()...)
		} else {
			result.Data.Datasets.Nodes = response.Nodes()
		}
		result.Data.Datasets.PageInfo.EndCursor = currentCursor
		resultNodeLen = len(result.Nodes())
		go writeResultFile(filename, &result)
		avgTime := time.Since(firstStart).Milliseconds() / int64(resultNodeLen-resultNodeStartLen)
		eta := avgTime * int64(response.NodeCount()-resultNodeLen)
		logStr := fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05.000 MST"), "]", currentCursor, resultNodeLen, response.NodeCount(), reqTime, time.Since(reqStartTime)-reqTime, time.Duration(avgTime)*time.Millisecond, time.Duration(eta)*time.Millisecond)
		//noinspection GoUnhandledErrorResult
		go logFile.Write([]byte(logStr))
		print(logStr)
	}
}

func writeResultFile(filename string, result *datacite.ResponseData) {
	resultString, err := json.Marshal(result)
	if err != nil {
		logError("Could not marshal JSON", err)
		panic(err)
	}
	if err = os.WriteFile(filename, resultString, 0644); err != nil {
		logError("Could not write File", err)
	}
}

func readResultFile(filename string, result *datacite.ResponseData) error {
	if filename == "" {
		return nil
	}
	fileData, err := os.ReadFile(filename)
	if err == nil {
		if err = json.Unmarshal(fileData, &result); err != nil {
			logError("File is not a JSON file", err)
			return err
		}
	} else {
		logError("Could not read any result file", err)
		return err
	}
	return nil
}

func logError(msg string, err error) {
	logStr := fmt.Sprintln("[", time.Now().Format("2006-01-02 15:04:05.000 MST"), "]", "ERROR", msg, ":", err.Error())
	println(logStr)
	_, err = logFile.Write([]byte(logStr))
	if err != nil {
		println("Error while writing log file", err.Error())
	}
}
