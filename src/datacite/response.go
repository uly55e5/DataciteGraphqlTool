package datacite

import "time"

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type ToolData struct {
	Query      string
	DataFormat string
	Time       time.Time
}

type ResponseData struct {
	Data struct {
		Datasets struct {
			TotalCount int         `json:"totalCount"`
			Nodes      interface{} `json:"nodes"`
			PageInfo   PageInfo    `json:"pageInfo"`
		} `json:"datasets"`
	} `json:"data"`
	Extensions map[string]interface{} `json:"extensions"`
	ToolData   ToolData
}

func (data ResponseData) Nodes() []interface{} {
	return data.Data.Datasets.Nodes.([]interface{})
}

func (data ResponseData) NodeCount() int {
	return data.Data.Datasets.TotalCount
}

func (data ResponseData) PageInfo() PageInfo {
	return data.Data.Datasets.PageInfo
}

type FullNode struct {
	Doi             string `json:"doi"`
	PublicationYear int    `json:"publicationYear"`
	Repository      struct {
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
}

func newFullResponse() ResponseData {
	data := ResponseData{}
	data.Data.Datasets.Nodes = []FullNode{}
	return data
}

type DoiNode struct {
	Doi string `json:"doi"`
}

func newDoiResponse() ResponseData {
	data := ResponseData{}
	data.Data.Datasets.Nodes = []DoiNode{}
	return data
}

var NamedResponse = map[string]ResponseData{
	"FullData": newFullResponse(),
	"DoiData":  newDoiResponse(),
}
