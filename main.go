package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type responseData struct {
	Data struct {
		Datasets struct {
			TotalCount int `json:"totalCount"`
			Nodes      []struct {
				Doi        string `json:"doi"`
				Repository struct {
					Name string `json:"name"`
				} `json:"repository"`
				References struct {
					Nodes []struct {
						Doi       string `json:"doi"`
						Container struct {
							Title string
						}
					} `json:"nodes"`
				} `json:"references"`
			} `json:"nodes"`
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				StartCursor string `json:"startCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
		} `json:"datasets"`
	} `json:"data"`
	Extensions map[string]interface{} `json:"extensions"`
}

type dataciteRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func main() {

	var reqData = dataciteRequest{
		Query: `query AllDatasets($cursorId: String){
  datasets(query: "chemistry",first:20, after: $cursorId) {
    totalCount
    nodes {
      doi
      repository {
        name
      }
      references(resourceTypeId: "text") {
        nodes {
          doi
          container {
            title
          }
        }
      }
    }
    pageInfo {
    endCursor
    startCursor
    hasNextPage
    }
  }
}`,
		Variables: map[string]string{"cursorId": ""},
	}
	result := responseData{}
	startLen := 0
	var startCursor string = ""
	filedata, err := ioutil.ReadFile("result.json")
	if err == nil {
		err = json.Unmarshal(filedata, &result)
		if err != nil {
			println(err.Error())
		} else {
			startCursor = result.Data.Datasets.PageInfo.EndCursor
			startLen = len(result.Data.Datasets.Nodes)
		}
	}
	firstStart := time.Now()
	for {
		start := time.Now()
		reqData.Variables["cursorId"] = startCursor
		reqValue, err := json.Marshal(reqData)
		if err != nil {
			println(err.Error())
			panic(err)
		}
		resp, err := http.Post("https://api.datacite.org/graphql", "application/json", bytes.NewBuffer(reqValue))
		if err != nil {
			println(err.Error())
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
			continue
		}
		jsonBody := responseData{}
		if err = json.Unmarshal(body, &jsonBody); err != nil {
			println(err.Error())
			continue
		}

		startCursor = jsonBody.Data.Datasets.PageInfo.EndCursor
		result.Data.Datasets.Nodes = append(result.Data.Datasets.Nodes, jsonBody.Data.Datasets.Nodes...)
		result.Data.Datasets.PageInfo.EndCursor = startCursor
		if !jsonBody.Data.Datasets.PageInfo.HasNextPage {
			break
		}
		resultString, err := json.Marshal(result)
		if err != nil {
			println(err.Error())
			panic(err)
		}
		err = ioutil.WriteFile("result.json", resultString, 0644)
		if err != nil {
			println(err.Error())
		}
		avgTime := time.Since(firstStart).Milliseconds() / int64(len(result.Data.Datasets.Nodes)-startLen)
		eta := avgTime * int64(jsonBody.Data.Datasets.TotalCount-len(result.Data.Datasets.Nodes))
		fmt.Println(startCursor, len(result.Data.Datasets.Nodes), jsonBody.Data.Datasets.TotalCount, time.Since(start), time.Duration(avgTime)*time.Millisecond, time.Duration(eta)*time.Millisecond)
	}
}
