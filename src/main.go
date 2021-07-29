package main

import (
	"encoding/json"
	"fmt"
	"github.com/uly55e5/DataciteDoiDownloader/datacite"
	"github.com/uly55e5/DataciteDoiDownloader/graphql"
	"io/ioutil"
	"time"
)

const filename = "result.json"

func main() {
	result := datacite.ResponseData{}
	startLen := 0
	var startCursor = ""

	err := readResultFile(filename, &result)
	if err == nil {
		startCursor = result.Data.Datasets.PageInfo.EndCursor
		startLen = len(result.Data.Datasets.Nodes)
	}
	firstStart := time.Now()
	reqData := datacite.FullDataRequest
	reqData.Variables["query"] = "chemistry"
	reqData.Variables["first"] = 20
	for {
		start := time.Now()
		response := datacite.ResponseData{}
		reqData.Variables["cursorId"] = startCursor
		if err := graphql.GetDataciteGraphQLResult(&reqData, &response); err != nil {
			continue
		}
		if !response.Data.Datasets.PageInfo.HasNextPage || response.Data.Datasets.PageInfo.EndCursor == "" {
			break
		}
		startCursor = response.Data.Datasets.PageInfo.EndCursor
		result.Data.Datasets.Nodes = append(result.Data.Datasets.Nodes, response.Data.Datasets.Nodes...)
		result.Data.Datasets.PageInfo.EndCursor = startCursor
		writeResultFile(filename, &result)
		avgTime := time.Since(firstStart).Milliseconds() / int64(len(result.Data.Datasets.Nodes)-startLen)
		eta := avgTime * int64(response.Data.Datasets.TotalCount-len(result.Data.Datasets.Nodes))
		fmt.Println(startCursor, len(result.Data.Datasets.Nodes), response.Data.Datasets.TotalCount, time.Since(start), time.Duration(avgTime)*time.Millisecond, time.Duration(eta)*time.Millisecond)
	}
}

func writeResultFile(filename string, result *datacite.ResponseData) {
	resultString, err := json.Marshal(result)
	if err != nil {
		println(err.Error())
		panic(err)
	}
	err = ioutil.WriteFile(filename, resultString, 0644)
	if err != nil {
		println(err.Error())
	}
}

func readResultFile(filename string, result *datacite.ResponseData) error {
	fileData, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(fileData, &result)
		if err != nil {
			println(err.Error())
			return err
		}
	} else {
		return err
	}
	return nil
}
